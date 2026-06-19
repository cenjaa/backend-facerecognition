import cv2
import os
import sys
import argparse
from pathlib import Path

def augment_image(img):
    """
    Applies random augmentations (flip, rotate, brightness) to an image.
    Returns a list of augmented images.
    """
    img_size = (100, 100)
    augs = []
    
    # 1. Original
    augs.append(img)
    
    # 2. Horizontal Flip (Simulates opposite side)
    augs.append(cv2.flip(img, 1))
    
    # 3. Random Rotations (Simulates head tilt)
    center = (img_size[0] // 2, img_size[1] // 2)
    for angle in [-15, 15]:
        M = cv2.getRotationMatrix2D(center, angle, 1.0)
        augs.append(cv2.warpAffine(img, M, img_size))
        
    # 4. Brightness Adjustments
    augs.append(cv2.convertScaleAbs(img, alpha=1.2, beta=10)) # Brighter
    augs.append(cv2.convertScaleAbs(img, alpha=0.8, beta=-10)) # Darker
    
    return augs

def extract_faces_from_video(video_path, user_id, output_dir="dataset", interval=5, max_images=100):
    """
    Extracts face images from a video file and saves them to a structured dataset folder.
    """
    video_path = Path(video_path)
    if not video_path.exists():
        print(f"Error: Video file not found at {video_path}")
        return

    user_dir = Path(output_dir) / str(user_id)
    user_dir.mkdir(parents=True, exist_ok=True)

    face_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + "haarcascade_frontalface_alt2.xml")
    cap = cv2.VideoCapture(str(video_path))
    if not cap.isOpened():
        print(f"Error: Could not open video file {video_path}")
        return

    print(f"Processing video: {video_path.name} for User ID: {user_id}")
    print(f"Saving to: {user_dir.absolute()}")

    # --- MANUAL ORIENTATION CHECK ---
    ret, first_frame = cap.read()
    if not ret:
        print("Error: Could not read video.")
        return

    preview = first_frame.copy()
    h, w = preview.shape[:2]
    if w > 800:
        s = 800 / w
        preview = cv2.resize(preview, (800, int(h * s)))

    print("\n==============================================")
    print("👉 A preview window has opened.")
    print("👉 Press 'r' in the window to rotate the video until it's upright.")
    print("👉 Press 'ENTER' or 'SPACE' to confirm and start cropping.")
    print("==============================================\n")

    rotations = [None, cv2.ROTATE_90_CLOCKWISE, cv2.ROTATE_180, cv2.ROTATE_90_COUNTERCLOCKWISE]
    rot_idx = 0

    while True:
        display = preview.copy()
        if rotations[rot_idx] is not None:
            display = cv2.rotate(display, rotations[rot_idx])
            
        cv2.putText(display, "Press 'r' to rotate. Press ENTER to confirm.", (10, 30), 
                    cv2.FONT_HERSHEY_SIMPLEX, 0.6, (0, 255, 0), 2)
        cv2.imshow("Orientation Check", display)
        
        key = cv2.waitKey(0) & 0xFF
        if key == ord('r'):
            rot_idx = (rot_idx + 1) % 4
        elif key == 13 or key == 32:  # Enter or Space
            break
            
    cv2.destroyAllWindows()
    locked_rotation = rotations[rot_idx]
    
    # Reset video back to frame 0
    cap.set(cv2.CAP_PROP_POS_FRAMES, 0)
    # --------------------------------

    count = 0
    frame_idx = 0
    img_size = (100, 100)

    while cap.isOpened() and count < max_images:
        ret, frame = cap.read()
        if not ret:
            break

        if frame_idx % interval == 0:
            # 1) Optimization: Scale down huge frames for detection
            h, w = frame.shape[:2]
            detect_frame = frame
            scale = 1.0
            if w > 800:
                scale = 800 / w
                detect_frame = cv2.resize(frame, (800, int(h * scale)))

            # 2) Apply the locked rotation
            if locked_rotation is not None:
                detect_frame = cv2.rotate(detect_frame, locked_rotation)

            gray = cv2.cvtColor(detect_frame, cv2.COLOR_BGR2GRAY)
            faces = face_cascade.detectMultiScale(gray, 1.3, 5, minSize=(50, 50))

            if len(faces) > 0:
                x, y, w, h = max(faces, key=lambda r: r[2] * r[3])
                face_img = cv2.resize(gray[y:y+h, x:x+w], img_size)
                
                # Apply Augmentations
                augmented_images = augment_image(face_img)
                
                for aug_img in augmented_images:
                    if count >= max_images:
                        break
                    output_filename = f"face_{count:04d}.jpg"
                    cv2.imwrite(str(user_dir / output_filename), aug_img)
                    count += 1
                
                if count % 20 == 0:
                    print(f"  -> Extracted {count} images (including augmentations)...")

        frame_idx += 1

    cap.release()
    print(f"\nDone! Extracted {count} images for User {user_id}.")
    print(f"Dataset path: {user_dir.absolute()}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Convert video to face dataset images.")
    parser.add_argument("--video", type=str, required=True, help="Path to the input video file")
    parser.add_argument("--id", type=int, required=True, help="User ID for the employee")
    parser.add_argument("--output", type=str, default="dataset", help="Root directory for dataset (default: 'dataset')")
    parser.add_argument("--interval", type=int, default=5, help="Extract every Nth frame (default: 5)")
    parser.add_argument("--max", type=int, default=100, help="Maximum number of images to extract (default: 100)")

    args = parser.parse_args()

    extract_faces_from_video(
        args.video, 
        args.id, 
        output_dir=args.output, 
        interval=args.interval, 
        max_images=args.max
    )
