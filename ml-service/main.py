import os
import shutil
import threading
import numpy as np
import cv2
import joblib
import yaml
from typing import Optional
from sklearn.svm import SVC
from sklearn.decomposition import PCA
from sklearn.model_selection import train_test_split
from sklearn.metrics import accuracy_score
from fastapi import FastAPI, UploadFile, File
from pydantic import BaseModel

import boto3
from botocore.config import Config as BotoConfig
from botocore.exceptions import ClientError

try:
    import pygad
    HAS_PYGAD = True
except ImportError:
    HAS_PYGAD = False

try:
    from rpca_algorithm import R_PCA
    HAS_RPCA = True
except ImportError:
    HAS_RPCA = False

BASE_DIR = os.path.abspath(os.path.dirname(__file__))
MODELS_DIR = os.path.join(BASE_DIR, "models")
DATASET_DIR = os.path.join(BASE_DIR, "dataset")

os.makedirs(MODELS_DIR, exist_ok=True)
os.makedirs(DATASET_DIR, exist_ok=True)


def load_config() -> dict:
    if os.environ.get("MINIO_HOST"):
        return {
            "minio": {
                "endpoint": f"{os.environ.get('MINIO_HOST', 'localhost')}:{os.environ.get('MINIO_PORT', '9000')}",
                "access_key": os.environ.get("MINIO_ACCESS_KEY", "admin"),
                "secret_key": os.environ.get("MINIO_SECRET_KEY", "Hsjdnvftrmm630!"),
                "bucket_name": os.environ.get("MINIO_BUCKET", "attendance"),
                "use_ssl": os.environ.get("MINIO_USE_SSL", "false").lower() == "true",
            }
        }

    config_path = os.path.join(BASE_DIR, "..", "config.yaml")
    if os.path.exists(config_path):
        with open(config_path, "r") as f:
            return yaml.safe_load(f)

    print("⚠️ No configuration found, using defaults")
    return {
        "minio": {
            "endpoint": "localhost:9000",
            "access_key": "admin",
            "secret_key": "Hsjdnvftrmm630!",
            "bucket_name": "attendance",
            "use_ssl": False,
        }
    }


config = load_config()


class S3Client:

    def __init__(self, minio_config: dict):
        self.bucket_name = minio_config.get("bucket_name", "attendance")

        endpoint = minio_config.get("endpoint", "localhost:9000")
        use_ssl = minio_config.get("use_ssl", False)
        protocol = "https" if use_ssl else "http"
        endpoint_url = f"{protocol}://{endpoint}"

        boto_config = BotoConfig(
            connect_timeout=5,
            read_timeout=30,
            retries={"max_attempts": 2},
        )

        self.s3 = boto3.client(
            "s3",
            endpoint_url=endpoint_url,
            aws_access_key_id=minio_config.get("access_key", ""),
            aws_secret_access_key=minio_config.get("secret_key", ""),
            config=boto_config,
        )

        self._ensure_bucket()

    def _ensure_bucket(self):
        try:
            self.s3.head_bucket(Bucket=self.bucket_name)
        except ClientError:
            try:
                self.s3.create_bucket(Bucket=self.bucket_name)
                print(f"☁️ Created bucket: '{self.bucket_name}'")
            except Exception as e:
                print(f"⚠️ Could not create bucket: {e}")

    def download_dataset(self, local_dir: str) -> int:
        os.makedirs(local_dir, exist_ok=True)
        paginator = self.s3.get_paginator("list_objects_v2")
        count = 0
        for page in paginator.paginate(Bucket=self.bucket_name, Prefix="dataset/"):
            for obj in page.get("Contents", []):
                key = obj["Key"]
                rel_path = key[len("dataset/"):]
                if not rel_path:
                    continue
                local_path = os.path.join(local_dir, rel_path.replace("/", os.sep))
                os.makedirs(os.path.dirname(local_path), exist_ok=True)
                self.s3.download_file(self.bucket_name, key, local_path)
                count += 1
        return count

    def upload_models(self, local_dir: str):
        for f in ["svm_model.pkl", "pca_transformer.pkl", "label_map.pkl"]:
            local_path = os.path.join(local_dir, f)
            if os.path.exists(local_path):
                self.s3.upload_file(local_path, self.bucket_name, f"models/{f}")
                print(f"   -> Uploaded {f}")


IMG_SIZE = (100, 100)

training_state = {
    "is_training": False,
    "progress": 0,
    "status": "idle",
    "message": "",
    "error": None,
}
_lock = threading.Lock()


def _set_state(status: str, message: str, progress: int = 0, error: str = None):
    with _lock:
        training_state["status"] = status
        training_state["message"] = message
        training_state["progress"] = progress
        training_state["error"] = error


def _augment(img_vec: np.ndarray) -> list:
    img = img_vec.reshape(IMG_SIZE)
    augmented = []

    # 1. Horizontal flip
    augmented.append(cv2.flip(img, 1).flatten())
    # 2. Brightness +20%
    augmented.append(np.clip(img.astype(np.float32) * 1.2, 0, 255).astype(np.uint8).flatten())
    # 3. Brightness -20%
    augmented.append(np.clip(img.astype(np.float32) * 0.8, 0, 255).astype(np.uint8).flatten())
    # 4. Rotation +10°
    center = (IMG_SIZE[0] // 2, IMG_SIZE[1] // 2)
    M_pos = cv2.getRotationMatrix2D(center, 10, 1.0)
    augmented.append(cv2.warpAffine(img, M_pos, IMG_SIZE).flatten())
    # 5. Rotation -10°
    M_neg = cv2.getRotationMatrix2D(center, -10, 1.0)
    augmented.append(cv2.warpAffine(img, M_neg, IMG_SIZE).flatten())
    # 6. Flipped + brightness
    flipped = cv2.flip(img, 1)
    augmented.append(np.clip(flipped.astype(np.float32) * 1.2, 0, 255).astype(np.uint8).flatten())

    return augmented


def _train(s3_client: Optional[S3Client]):
    """Full training pipeline — runs in a background thread."""
    try:
        with _lock:
            training_state["is_training"] = True

        # Step 1: Download dataset
        _set_state("training", "Downloading dataset from S3...", 5)
        if s3_client:
            count = s3_client.download_dataset(DATASET_DIR)
            print(f"📥 Downloaded {count} files from S3")

        # Step 2: Load images
        _set_state("training", "Loading images...", 15)
        data, labels = [], []
        label_map = {}

        if not os.path.exists(DATASET_DIR):
            raise Exception("No dataset directory found")

        folders = [
            f for f in os.listdir(DATASET_DIR)
            if os.path.isdir(os.path.join(DATASET_DIR, f))
        ]
        if not folders:
            raise Exception("No user folders found in dataset")

        for folder_id in folders:
            try:
                uid = int(folder_id)
            except ValueError:
                continue

            label_map[uid] = uid
            path = os.path.join(DATASET_DIR, folder_id)

            for img_name in os.listdir(path):
                try:
                    img = cv2.imread(os.path.join(path, img_name), 0)
                    if img is None:
                        continue
                    vec = cv2.resize(img, IMG_SIZE).flatten()
                    data.append(vec)
                    labels.append(uid)
                except Exception:
                    pass

        if len(data) < 2:
            raise Exception(f"Not enough training data (found {len(data)} images)")

        # Step 3: Augmentation
        _set_state("training", "Augmenting dataset...", 25)
        original_count = len(data)
        augmented_data, augmented_labels = [], []
        for vec, lbl in zip(data, labels):
            for aug_vec in _augment(vec):
                augmented_data.append(aug_vec)
                augmented_labels.append(lbl)
        data.extend(augmented_data)
        labels.extend(augmented_labels)

        X = np.array(data)
        y = np.array(labels)
        print(f"📊 Dataset: {original_count} originals → {X.shape[0]} after augmentation, {len(set(y))} users")

        # Step 4: RPCA denoising
        if HAS_RPCA:
            _set_state("training", f"Running RPCA on {X.shape} matrix...", 35)
            rpca = R_PCA(X.T)
            L, S = rpca.fit(max_iter=100)
            X_clean = L.T
        else:
            print("⚠️ RPCA not available, skipping denoising")
            X_clean = X

        # Step 5: PCA dimensionality reduction
        _set_state("training", "Dimensionality reduction (PCA)...", 50)
        n_components = min(100, X_clean.shape[0] - 1, X_clean.shape[1])
        pca = PCA(n_components=n_components, whiten=True)
        X_pca = pca.fit_transform(X_clean)

        # Step 6: GA-SVM optimization
        X_train, X_test, y_train, y_test = train_test_split(
            X_pca, y, test_size=0.2, random_state=42
        )

        if HAS_PYGAD and len(X_train) >= 5:
            _set_state("training", "Optimizing SVM with Genetic Algorithm...", 60)

            def fitness_func(ga_instance, solution, solution_idx):
                C, gamma = solution[0], solution[1]
                clf = SVC(C=C, gamma=gamma, kernel="rbf")
                clf.fit(X_train, y_train)
                return accuracy_score(y_test, clf.predict(X_test))

            ga_instance = pygad.GA(
                num_generations=20,
                num_parents_mating=5,
                fitness_func=fitness_func,
                sol_per_pop=15,
                num_genes=2,
                gene_space=[
                    {"low": 0.1, "high": 100},
                    {"low": 0.0001, "high": 1},
                ],
            )
            ga_instance.run()
            best_solution, best_fitness, _ = ga_instance.best_solution()
            best_C, best_gamma = best_solution[0], best_solution[1]
            print(f"   -> Best Params: C={best_C:.2f}, Gamma={best_gamma:.4f}")
        else:
            _set_state("training", "Training SVM...", 60)
            best_C, best_gamma = 10.0, 0.01

        # Step 7: Train final model
        _set_state("training", "Training final model...", 75)
        final_svm = SVC(C=best_C, gamma=best_gamma, kernel="rbf", probability=True)
        final_svm.fit(X_pca, y)

        # Step 8: Save models
        _set_state("training", "Saving models...", 85)
        os.makedirs(MODELS_DIR, exist_ok=True)
        joblib.dump(final_svm, os.path.join(MODELS_DIR, "svm_model.pkl"))
        joblib.dump(pca, os.path.join(MODELS_DIR, "pca_transformer.pkl"))
        joblib.dump(label_map, os.path.join(MODELS_DIR, "label_map.pkl"))

        # Step 9: Upload to S3
        if s3_client:
            _set_state("training", "Uploading models to S3...", 90)
            s3_client.upload_models(MODELS_DIR)

        # Step 10: Clean up local dataset
        try:
            shutil.rmtree(DATASET_DIR)
            print("🗑️ Deleted local dataset/ after training")
        except Exception as e:
            print(f"⚠️ Could not delete local dataset: {e}")

        # Calculate accuracy
        accuracy = 0.0
        if len(X_test) > 0:
            try:
                accuracy = accuracy_score(
                    y_test,
                    final_svm.predict(pca.transform(X_clean[len(X_train):len(X_train) + len(X_test)])),
                )
            except Exception:
                pass

        _set_state("completed", f"Training complete! Accuracy: {accuracy:.1%}", 100)
        print(f"✅ Training Complete! Accuracy: {accuracy:.1%}")
        
        # Reload models into memory
        load_models()

    except Exception as e:
        _set_state("error", str(e), 0, error=str(e))
        print(f"❌ Training failed: {e}")
        import traceback
        traceback.print_exc()
    finally:
        with _lock:
            training_state["is_training"] = False


# ─── FastAPI App ───────────────────────────────────────
app = FastAPI(title="Face Recognition ML Service", version="1.0.0")

# Initialize S3 client
try:
    s3_client = S3Client(config.get("minio", {}))
    print("☁️ S3/MinIO client initialized")
except Exception as e:
    print(f"⚠️ S3/MinIO unavailable: {e}")
    s3_client = None

# Global Model Variables
global_svm = None
global_pca = None
global_label_map = None
face_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + "haarcascade_frontalface_alt2.xml")

def load_models():
    """Load models from disk into memory"""
    global global_svm, global_pca, global_label_map
    try:
        if os.path.exists(os.path.join(MODELS_DIR, "svm_model.pkl")):
            global_svm = joblib.load(os.path.join(MODELS_DIR, "svm_model.pkl"))
            global_pca = joblib.load(os.path.join(MODELS_DIR, "pca_transformer.pkl"))
            global_label_map = joblib.load(os.path.join(MODELS_DIR, "label_map.pkl"))
            print("🧠 Models loaded into memory for inference")
        else:
            print("⚠️ No models found locally. Waiting for sync or training.")
    except Exception as e:
        print(f"⚠️ Failed to load models: {e}")

# Load models on startup
load_models()


class TrainResponse(BaseModel):
    status: str
    message: str


class StatusResponse(BaseModel):
    status: str
    progress: int
    message: str
    error: Optional[str] = None


@app.get("/health")
def health():
    return {"status": "ok", "service": "ml-training"}

@app.post("/infer")
async def infer_face(file: UploadFile = File(...)):
    """Run face detection and recognition on uploaded image."""
    global global_svm, global_pca, global_label_map
    
    if global_svm is None or global_pca is None:
        return {"status": "error", "message": "Models not loaded"}
        
    try:
        # Read image
        contents = await file.read()
        nparr = np.frombuffer(contents, np.uint8)
        frame = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
        
        if frame is None:
            return {"status": "error", "message": "Invalid image"}
            
        # Convert to grayscale for Haar Cascade
        gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
        
        # Save frame to disk for debugging
        cv2.imwrite("debug_last_frame.jpg", frame)
        
        # Detect faces (using parameters optimized for speed/accuracy)
        faces = face_cascade.detectMultiScale(gray, 1.1, 5, minSize=(30, 30))
        
        if len(faces) == 0:
            print(f"👁️ /infer - No face detected in frame {frame.shape}")
            return {"status": "success", "detected": False, "message": "No face detected"}
            
        print(f"👁️ /infer - Detected {len(faces)} face(s)")
            
        # Find the largest face by bounding box area
        largest_face = max(faces, key=lambda rect: rect[2] * rect[3])
        x, y, w, h = largest_face
        
        # Extract face ROI and resize to expected model input
        face_roi = cv2.resize(gray[y:y+h, x:x+w], IMG_SIZE)
        
        # Transform and predict
        processed = global_pca.transform(face_roi.reshape(1, -1))
        probs = global_svm.predict_proba(processed)[0]
        
        max_prob = float(np.max(probs))
        best_idx = np.argmax(probs)
        
        predicted_label = int(global_svm.classes_[best_idx])
        user_id = global_label_map.get(predicted_label, predicted_label)
        confidence_pct = max_prob * 100
        
        print(f"   -> Match: User {user_id} with {confidence_pct:.1f}% confidence")
        
        return {
            "status": "success",
            "detected": True,
            "user_id": user_id,
            "confidence": max_prob * 100, # return as percentage
            "box": [int(x), int(y), int(w), int(h)]
        }
        
    except Exception as e:
        import traceback
        traceback.print_exc()
        return {"status": "error", "message": str(e)}


@app.post("/train", response_model=TrainResponse)
def start_training():
    """Trigger model training in a background thread."""
    with _lock:
        if training_state["is_training"]:
            return TrainResponse(
                status="error",
                message="Training already in progress",
            )

    thread = threading.Thread(target=_train, args=(s3_client,), daemon=True)
    thread.start()

    return TrainResponse(
        status="started",
        message="Training started in background",
    )


@app.get("/train_status", response_model=StatusResponse)
def get_train_status():
    """Poll training progress."""
    with _lock:
        return StatusResponse(
            status=training_state["status"],
            progress=training_state["progress"],
            message=training_state["message"],
            error=training_state["error"],
        )


if __name__ == "__main__":
    import uvicorn

    port = int(os.environ.get("ML_SERVICE_PORT", "8001"))
    print(f"🚀 ML Service starting on port {port}")
    uvicorn.run(app, host="0.0.0.0", port=port)
