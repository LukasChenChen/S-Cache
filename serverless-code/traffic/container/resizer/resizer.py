from PIL import Image
import logging

def resize():
    # Opens a image in RGB mode
    im = Image.open("test.jpg")
 
    # Size of the image in pixels (size of original image)
    # (This is not mandatory)
    width, height = im.size
 
    # Setting the points for cropped image
    left = 4
    top = height / 5
    right = 154
    bottom = 3 * height / 5
 
    # Cropped image of above dimension
    # (It will not change original image)
    im1 = im.crop((left, top, right, bottom))
    newsize = (300, 300)
    im1 = im1.resize(newsize)

    # save a image using extension
    im1 = im1.save("test-resized.jpg")
    logging.warning('resize job finish')

if __name__ == "__main__":
   resize()