#Checkout program
import logging

def checkout(barcode):
    items = {'banana':{'price': 4,   'stock': 6, 'code':123 },
         'apple':{'price': 2,   'stock': 0,'code':1231 },
         'orange':{'price': 1.5, 'stock': 32,'code':1233},
         'pear':{'price': 3,   'stock': 15,'code':12335},}

    bar = barcode
    sum = 0
    price =[]
    products =[]
    for key in items:
        if items[key]['code'] == bar:
            # print("item found")
            # print (key)
            # print ("price: %s" % items[key]['price'])
            # print ("stock: %s" % items[key]['stock'])
            # print ("Barcode: %s" % items[key]['code'])
            price.append(items[key]['price'])
            products.append(key)

            logging.warning('item found, price : %d',items[key]['price'])
            
            logging.warning("Finish Checkout....................")

if __name__ == "__main__":
   checkout(12335)