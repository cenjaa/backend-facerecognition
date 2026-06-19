import cv2
import os
import sys
import argparse
from pathlib import Path

def manual_video_to_dataset(video_path, user_id, output_dir="dataset"):
    """
    Manually extract poses from a video file.
    Use trackbar to scrub, 'C' to capture 10 frames, 'N' for next pose.
    """
    user_dir = Path(output_dir) / str(user_id)
    user_dir.mkdir(parents=True, exist_ok=True)

    cap = cv2.VideoCapture(video_path)
    if not cap.isOpened():
        print(f"Error: Could not open video {video_path}")
        return

    # --- MANUAL ORIENTATION CHECK ---
    ret, first_frame = cap.read()
    if not ret:
        print("Error: Could not read video.")
        return

    preview = first_frame.copy()
    ph, pw = preview.shape[:2]
    if pw > 800:
        ps = 800 / pw
        preview = cv2.resize(preview, (800, int(ph * ps)))

    print("\n==============================================")
    print("👉 A preview window has opened.")
    print("👉 Press 'r' in the window to rotate the video until it's upright.")
    print("👉 Press 'ENTER' or 'SPACE' to confirm and start scrubbing.")
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

    total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
    face_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + "haarcascade_frontalface_alt2.xml")
    
    poses = ["FRONT", "LEFT", "RIGHT", "UP", "DOWN"]
    current_pose_idx = 0
    img_size = (100, 100)
    total_captured = 0
    
    paused = True
    frame_pos = 0

    def on_trackbar(val):
        nonlocal frame_pos
        frame_pos = val
        cap.set(cv2.CAP_PROP_POS_FRAMES, frame_pos)

    window_name = "Manual Video Extraction"
    cv2.namedWindow(window_name)
    cv2.createTrackbar("Position", window_name, 0, total_frames - 1, on_trackbar)

    print(f"--- Manual Extraction for User ID: {user_id} ---")
    print("Controls:")
    print(" [SPACE] : Play/Pause")
    print(" [C]     : Capture 10 frames from current position")
    print(" [N]     : Switch to Next Pose")
    print(" [Q]     : Quit")
    print("------------------------------------------")

    while True:
        if not paused:
            frame_pos = int(cap.get(cv2.CAP_PROP_POS_FRAMES))
            cv2.setTrackbarPos("Position", window_name, frame_pos)
        
        ret, frame = cap.read()
        if not ret:
            paused = True
            cap.set(cv2.CAP_PROP_POS_FRAMES, total_frames - 1)
            continue

        # Apply Rotation and Resizing for UI
        if locked_rotation is not None:
            frame = cv2.rotate(frame, locked_rotation)
        
        h, w = frame.shape[:2]
        if w > 1000: # Scale down if too large for screen
            s = 1000 / w
            frame = cv2.resize(frame, (1000, int(h * s)))

        display_frame = frame.copy()
        gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
        
        # UI Overlay
        pose_label = poses[current_pose_idx] if current_pose_idx < len(poses) else "ALL DONE"
        status = "PAUSED" if paused else "PLAYING"
        cv2.putText(display_frame, f"POSE: {pose_label} | {status}", (20, 40), cv2.FONT_HERSHEY_DUPLEX, 0.8, (0, 255, 0), 2)
        cv2.putText(display_frame, f"Captured: {total_captured}/50", (20, 80), cv2.FONT_HERSHEY_DUPLEX, 0.7, (255, 255, 255), 1)

        cv2.imshow(window_name, display_frame)
        
        key = cv2.waitKey(30) & 0xFF
        if key == ord(' '):
            paused = not paused
        elif key == ord('q'):
            break
        elif key == ord('n'):
            current_pose_idx = (current_pose_idx + 1) % len(poses)
            print(f"Switched to pose: {poses[current_pose_idx]}")
        elif key == ord('c'):
            print(f"Capturing 10 frames for {poses[current_pose_idx]}...")
            captured_now = 0
            # Capture 10 frames from this point
            temp_pos = frame_pos
            while captured_now < 10:
                success, f = cap.read()
                if not success: break
                
                # Apply rotation for capture
                if locked_rotation is not None:
                    f = cv2.rotate(f, locked_rotation)
                
                g = cv2.cvtColor(f, cv2.COLOR_BGR2GRAY)
                faces = face_cascade.detectMultiScale(g, 1.3, 5, minSize=(50, 50))
                
                if len(faces) > 0:
                    x, y, w, h = max(faces, key=lambda r: r[2] * r[3])
                    face_img = cv2.resize(g[y:y+h, x:x+w], img_size)
                    
                    # Use simple numeric naming (1, 2, 3...)
                    filename = f"{total_captured + 1}.jpg"
                    cv2.imwrite(str(user_dir / filename), face_img)
                    captured_now += 1
                    total_captured += 1
                
                # Skip 2 frames to get slight variation
                cap.set(cv2.CAP_PROP_POS_FRAMES, cap.get(cv2.CAP_PROP_POS_FRAMES) + 2)
            
            # Go back to where we were
            cap.set(cv2.CAP_PROP_POS_FRAMES, temp_pos)
            print(f"Total Captured: {total_captured}/50")

    cap.release()
    cv2.destroyAllWindows()
    print(f"\nDone! Saved {total_captured} images to {user_dir}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Manual Video to Dataset Tool")
    parser.add_argument("--video", type=str, required=True, help="Path to video file")
    parser.add_argument("--id", type=int, required=True, help="User ID")
    
    args = parser.parse_args()
    manual_video_to_dataset(args.video, args.id)
