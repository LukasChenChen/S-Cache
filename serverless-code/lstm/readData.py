#!/usr/bin/env python3
import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
import math
import csv

##extract the number of invocations for a specific function 
def read_azure_2021(file_name = 'AzureFunctionTrace.txt', func_id = '48cc770d590d3c5a7691b3b4e9302f82ec3be5ddc2a037d94ad2e76f44dd8946'):

    with open(file_name,'r') as f:
        contents = f.readlines()
    app = []
    func = []
    timestamp = []
    duration = []
    startTime = []
    roundStartTime = []
#the size of all time slot
    count  = [0]*1209600

    for index in range(1,len(contents)):
        parameters = contents[index].split(",")
        app.append(parameters[0])
        func.append(parameters[1])
        timestamp.append(float(parameters[2]))
        duration.append(float(parameters[3]))
        startTime.append(float(parameters[2])-float(parameters[3]))
        roundStartTime.append(int(math.floor(float(parameters[2]) - float(parameters[3]))))
        if parameters[1] == func_id:
            count[roundStartTime[-1]] += 1



    data = pd.DataFrame(np.array([count]))
    data = data.transpose()
    data.columns = ['count']

    return data

def read_azure_2019(file_name, function_id):
   
    with open(file_name,'r') as f:
        contents = f.readlines()
    owner = []
    app = []
    function = []
    trigger = []
    time_slot = []

    for index in range(1,len(contents)):
        parameters = contents[index].split(",")
        if parameters[2] == function_id:
        
            owner.append(parameters[0])
            app.append(parameters[1])
            function.append(parameters[2])
            trigger.append(parameters[3])
            time_slot = parameters[4:]
            break
        



    data = pd.DataFrame(np.array([time_slot], dtype = np.float32))
#test_data_size = 120960
    data = data.transpose()
    data.columns = ['Invocation']

    return data

def read_file(file_name = './request/final/sortedd01.csv'):
#     with open(file_name,'r') as f:
#         contents = f.readlines()
#     app = []
#     totalNum = []
#     time_slot = []
#     data = pd.DataFrame( dtype = np.float32)
# #the size of all time slot
#     count  = [0]*1209600
    

#     titles = contents[0]
#     print(titles)
#     for index in range(1,len(contents)):
#         parameters = contents[index].split(",")
#         print(parameters[0])
#         app.append(parameters[0])
#         totalNum.append(parameters[1])
#         time_slot = parameters[2:]
#         tmp = pd.DataFrame(np.array([time_slot], dtype = np.float32))
#         tmp = tmp.transpose()
#         pd.concat([tmp, data])

#     data.columns = titles

    df = pd.read_csv(file_name) 
    df = df.transpose()
    # print(df)

    return df

# sort the data based on frequence.
def consolidate_data(file_names, output_name):

    for file_name in file_names:
        with open(file_name,'r') as f:
            contents = f.readlines()

        data_map = {}
        time_slot = [1440]
        
        for index in range(1,len(contents)):
            parameters = contents[index].split(",")
            app_id = parameters[1]
            time = parameters[4:]
            #convert to int
            time_int = list(map(int, time))
            time_slot = np.array(time_int)
                
            if data_map.get(app_id) is None:
            #new key
                data_map.update({app_id : time_slot})
            
            else:
            #old key
                current_time_slot = data_map[app_id]
                update_time_slot = np.add(time_slot, current_time_slot)
                data_map.update({app_id : update_time_slot})


    ## write data_map to excel
    f = open(output_name, 'w')

    # create the csv writer
    writer = csv.writer(f)

    f.write('App')
    f.write(',')
    f.write('totalNumber')
    f.write(',')
    time_slot_index = list(np.arange(1,1440+1))
    writer.writerow(time_slot_index)

    for key, value in data_map.items():
        f.write(key)
        f.write(',')
        f.write("%s" % np.sum(value))
        f.write(',')
        writer.writerow(value)
    f.close()

def write_file(output_name, data_map):
     ## write data_map to excel
    f = open(output_name, 'w')

    # create the csv writer
    writer = csv.writer(f)

    f.write('App')
    f.write(',')
    f.write('totalNumber')
    f.write(',')
    time_slot_index = list(np.arange(1,1440+1))
    writer.writerow(time_slot_index)

    for key, value in data_map.items():
        f.write(key)
        f.write(',')
        f.write("%s" % np.sum(value))
        f.write(',')
        writer.writerow(value)
    f.close()

#read 13 files, merge the same appID, calculate the total number.
def consolidate_all_files(output_name):
    file_names =[]
    # 1- 13
    for i in range(1,14):
        if i < 10:
            file_index = '0'+str(i)
        else:
            file_index = str(i)
        file_name = './request/invocations_per_function_md.anon.d' + file_index + '.csv'
        print(file_name)
        file_names.append(file_name)
    
    consolidate_data(file_names, output_name)

# sort the files based on number of innovcations
def sort_data(file_name, output_name):
    # DataFrame to read our input CS file
    dataFrame = pd.read_csv(file_name)

    print(dataFrame.head(3))

    # sorting according to totalNumber column
    dataFrame = dataFrame.sort_values(by = ['totalNumber'], ascending=False, na_position = 'last')
   
    dataFrame.head(4).to_csv(output_name, index = False)
    
    print('after sorting')
    print(dataFrame.head(3))

# find apps with a particular id in all files, return data_map
def find_app(file_names, ids):
    
    file_count = 1
    for file_name in file_names:
        with open(file_name,'r') as f:
            contents = f.readlines()

        data_map = {}
        time_slot = [1440]
        
        for index in range(1,len(contents)):
            parameters = contents[index].split(",")
            app_id = parameters[1]
            if app_id not in ids:
                continue
            else:
                time = parameters[4:]
                #convert to int
                time_int = list(map(int, time))
                time_slot = np.array(time_int)
                
                if data_map.get(app_id) is None:
                    #new key
                    data_map.update({app_id : time_slot})
            
                else:
                    #old key
                    current_time_slot = data_map[app_id]
                    update_time_slot = np.add(time_slot, current_time_slot)
                    data_map.update({app_id : update_time_slot})
        
        write_file('./request/final/'+ 'invocations_per_function_md.anon.d'+ str(file_count) + '.csv', data_map)

        file_count += 1


#get the top k most invoked app_id from './request/sorted.csv'
def get_top_apps(file_name, k):
     df = pd.read_csv(file_name)
     app_ids = df.loc[:, 'App'].tolist()

     return app_ids[: k]

def find_apps():
    app_ids = get_top_apps('/request/final/sortedd01.csv', 4)

    file_names =[]
    # 1- 13
    for i in range(1,14):
        if i < 10:
            file_index = '0'+str(i)
        else:
            file_index = str(i)
        file_name = './request/invocations_per_function_md.anon.d' + file_index + '.csv'
        print(file_name)
        file_names.append(file_name)

    find_app(file_names, app_ids)

#main programm
#consolidate 13 files
#consolidate_all_files('./request/request_data_set.csv')
def create_top_4():
    filenames = []
    filename = './request/invocations_per_function_md.anon.d01.csv'
    filenames.append(filename)
    consolidate_data(filenames, './request/final/mergedd01.csv')
    #sort the files
    sort_data('./request/final/mergedd01.csv', './request/final/sortedd01.csv')

# find_apps()

#get the most used apps, output 13 files that only contain these apps

def main():
    #create_top_4()
    read_file()

if __name__ == "__main__":
   main()