from ultralytics import YOLO

def infer(image_path):
    # Load a pretrained YOLOv8n-pose Pose model
    model = YOLO("/test_mount/test-pv/yolov8n-cls.pt")

    results = model(image_path)  # results list
    print("---------")
    name =  model.names[results[0].probs.top1]
    return {"result": name}

def main(userparams):
    path = userparams.get("full_path")
    return infer(image_path=path)