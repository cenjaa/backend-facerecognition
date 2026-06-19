import cv2
import os
import sys
import time
import argparse
from pathlib import Path

def run_face_registration(user_id, output_dir="dataset"):
    """
    Interactive face registration tool.
    Captures 10 images for each of the 5 poses (Total 50 images).
    """
    # 1. Setup paths
    user_dir = Path(output_dir) / str(user_id)
    user_dir.mkdir(parents=True, exist_ok=True)

    # 2. Initialize Camera and Detector
    cap = cv2.VideoCapture(0)
    if not cap.isOpened():
        print("Error: Could not open camera.")
        return

    face_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + "haarcascade_frontalface_alt2.xml")
    
    poses = ["FRONT", "LEFT", "RIGHT", "UP", "DOWN"]
    images_per_pose = 10
    img_size = (100, 100)
    total_captured = 0

    print(f"--- Face Registration for User ID: {user_id} ---")
    print("Controls:")
    print(" [SPACE] : Start capturing for current pose")
    print(" [Q]     : Quit")
    print("------------------------------------------")

    for pose in poses:
        captured_for_pose = 0
        is_capturing = False
        
        while captured_for_pose < images_per_pose:
            ret, frame = cap.read()
            if not ret:
                break

            # Mirror for natural feel
            frame = cv2.flip(frame, 1)
            display_frame = frame.copy()
            gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
            
            # Detect face
            faces = face_cascade.detectMultiScale(gray, 1.3, 5, minSize=(100, 100))
            
            # Draw UI
            status_text = f"POSE: {pose} ({captured_for_pose}/{images_per_pose})"
            color = (0, 255, 0) if is_capturing else (0, 165, 255)
            
            cv2.putText(display_frame, status_text, (30, 50), cv2.FONT_HERSHEY_DUPLEX, 1, color, 2)
            if not is_capturing:
                cv2.putText(display_frame, "Press [SPACE] to start", (30, 90), cv2.FONT_HERSHEY_DUPLEX, 0.7, (255, 255, 255), 1)

            if len(faces) > 0:
                # Take largest face
                x, y, w, h = max(faces, key=lambda r: r[2] * r[3])
                cv2.rectangle(display_frame, (x, y), (x+w, y+h), color, 2)
                
                if is_capturing:
                    # Use simple numeric naming (1, 2, 3...)
                    face_img = cv2.resize(gray[y:y+h, x:x+w], img_size)
                    filename = f"{total_captured + 1}.jpg"
                    cv2.imwrite(str(user_dir / filename), face_img)
                    
                    captured_for_pose += 1
                    total_captured += 1
                    time.sleep(0.1) # Small delay to get slightly different angles
            
            cv2.imshow("PNM Face Registration", display_frame)
            
            key = cv2.waitKey(1) & 0xFF
            if key == ord(' '):
                is_capturing = True
            elif key == ord('q'):
                print("Registration cancelled.")
                cap.release()
                cv2.destroyAllWindows()
                return

        print(f"Done with pose: {pose}")
        is_capturing = False
        time.sleep(1) # Pause between poses

    cap.release()
    cv2.destroyAllWindows()
    print(f"\nSuccessfully captured 50 images for User {user_id}!")
    print(f"Dataset path: {user_dir.absolute()}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Interactive Face Registration Tool")
    parser.add_argument("--id", type=int, required=True, help="User ID for the employee")
    parser.add_argument("--output", type=str, default="dataset", help="Output directory")

    args = parser.parse_args()
    run_face_registration(args.id, args.output)
