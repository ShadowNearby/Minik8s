from PIL import Image

def resize_pic(image_name, width, height):
  input_image_path = '/test_mount/test-pv/'+ image_name
  image = Image.open(input_image_path)

  new_size = (width, height)
  image.thumbnail(new_size)

  output_image_path = '/test_mount/test-pv/'+image_name
  image.save(output_image_path)
  return {"full_path": output_image_path}


new_width = 320
new_height = 240

def main(userparams):
    image_path= userparams.get("path")
    return resize_pic(image_path, new_width, new_height)