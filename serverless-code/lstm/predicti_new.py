#using lstm to predict data from azure_2019 trace
#author: chen
#01.07.2022
import torch
import torch.nn as nn
import seaborn as sns
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
import time
from statistics import mean
from readData import read_azure_2019
from readData import read_file
import csv
#%matplotlib inline


def create_inout_sequences(input_data, tw):
    inout_seq = []
    L = len(input_data)
    for i in range(L-tw):
        train_seq = input_data[i:i+tw]
        train_label = input_data[i+tw:i+tw+1]
        inout_seq.append((train_seq ,train_label))
    return inout_seq



class LSTM(nn.Module):
    def __init__(self, input_size=1, hidden_layer_size=100, output_size=1, bias=False):
        super().__init__()
        self.hidden_layer_size = hidden_layer_size

        self.lstm = nn.LSTM(input_size, hidden_layer_size, bias=False)

        self.linear = nn.Linear(hidden_layer_size, output_size, bias=False)

        self.hidden_cell = (torch.zeros(1,1,self.hidden_layer_size),
                            torch.zeros(1,1,self.hidden_layer_size))

    def forward(self, input_seq):
        lstm_out, self.hidden_cell = self.lstm(input_seq.view(len(input_seq) ,1, -1), self.hidden_cell)
        predictions = self.linear(lstm_out.view(len(input_seq), -1))
        return predictions[-1]
def predict(in_data):
    # #file name and function id
    # file_name = './request/invocations_per_function_md.anon.d01.csv'
    # function_id = '5e98666f0af3fee4a98e670ab893ddf57046816b30775e3decd77a098d317e98'

    # in_data = read_azure_2019(file_name, function_id)


    fig_size = plt.rcParams["figure.figsize"]
    fig_size[0] = 15
    fig_size[1] = 5
    plt.rcParams["figure.figsize"] = fig_size




    all_data = in_data['Invocation'].values.astype(float)

    # data size is 1440 for each function
    #one data represent 1 minute
    #use 20% data as testset
    test_data_size = 288

    train_data = all_data[:-test_data_size]
    test_data = all_data[-test_data_size:]




    print(len(train_data))
    print(len(test_data))

    start_train_time = time.time()
    #normalize data to [-1,1]
    from sklearn.preprocessing import MinMaxScaler

    from sklearn.preprocessing import StandardScaler
    #scaler = MinMaxScaler(feature_range=(-1, 1))
    #train_data_normalized = scaler.fit_transform(train_data .reshape(-1, 1))

    scaler = StandardScaler()
    train_data_normalized = scaler.fit_transform(train_data.reshape(-1,1))



    train_data_normalized = torch.FloatTensor(train_data_normalized).view(-1)
    ##TODO:
    #set the window to be 10minute
    # train_window = 30
    train_window = 60

    train_inout_seq = create_inout_sequences(train_data_normalized, train_window)

    train_inout_seq[:5]

    model = LSTM()
    loss_function = nn.MSELoss()
    optimizer = torch.optim.Adam(model.parameters(), lr=0.001)

    print(model)

    epochs = 150

    for i in range(epochs):
        for seq, labels in train_inout_seq:
            optimizer.zero_grad()
            model.hidden_cell = (torch.zeros(1, 1, model.hidden_layer_size),
                            torch.zeros(1, 1, model.hidden_layer_size))

            y_pred = model(seq)

            single_loss = loss_function(y_pred, labels)
            single_loss.backward()
            optimizer.step()

        if i%2 == 1:
            print(f'epoch: {i:3} loss: {single_loss.item():10.8f}')

    #print(f'epoch: {i:3} loss: {single_loss.item():10.10f}')
    end_train_time = time.time()

    start_pred_time = time.time()

    #indicate how many values we predict
    fut_pred = test_data_size

    ##decide how many values we use as initial input
    #test_inputs = train_data_normalized[-train_window:].tolist()
    test_inputs = train_data_normalized[-train_window:].tolist()
    print('test_inputs')
    print(test_inputs)

    model.eval()

    #append the predict value every time
    for i in range(fut_pred):
        seq = torch.FloatTensor(test_inputs[-train_window:])
        with torch.no_grad():
            model.hidden = (torch.zeros(1, 1, model.hidden_layer_size),
                            torch.zeros(1, 1, model.hidden_layer_size))
            test_inputs.append(model(seq).item())

    print('size of test_inputs')
    print(len(test_inputs))
    print('predict value')
    print(test_inputs[-fut_pred:])

    ##reverse the value to real values
    actual_predictions = scaler.inverse_transform(np.array(test_inputs[-test_data_size:] ).reshape(-1, 1))
    #print('actual predictions')
    #print(actual_predictions)

    end_pred_time = time.time()

    print("train time is %s" % (end_train_time-start_train_time))
    print("prediction time is %s"% (end_pred_time - start_pred_time))

    x = np.arange(1152, 1440, 1)


    error_rate = [abs(i-j)/j*100 for i,j in zip(actual_predictions,test_data)]
    print('The average error is %s ' % (sum(error_rate) / len(error_rate)))

    with open('./request/final/errors.csv', "a") as file:
        writer = csv.writer(file, delimiter=',')
        writer.writerow((sum(error_rate) / len(error_rate)))
    
    actual_predict_df = pd.DataFrame(actual_predictions)
    actual_predict_df = actual_predict_df.transpose()
    actual_predict_df.to_csv('./request/final/prediction.csv', encoding='utf-8', index=False, mode='a')
    
    
    # plt.subplot(2,1,1)
    # plt.title('Time  vs Inovcation')
    # plt.ylabel('Invocations')
    # plt.grid(True)
    # plt.autoscale(axis='x', tight=True)
    # plt.plot(in_data['Invocation'], 'r')
    # plt.plot(x,actual_predictions, 'g')

    # plt.subplot(2,1,2)
    # plt.title('Error')
    # plt.xlabel('Time')
    # plt.ylabel('Error in percentage')
    # plt.plot(x,error_rate, 'b')
    # #plt.plot(actual_predictions, 'g')
    # plt.show()

## predict the trend for 4
def predict_all(filename):
    df = read_file(filename)
    print('length', len(df.columns))
    
    for column in df:
        # print(df[column][2:])
        tmp = df[column][2:]
        tmp_df = pd.DataFrame({'Invocation':tmp})
        # print(tmp_df)
        predict(tmp_df)
        # break
   

def main():
    #create_top_4()
    predict_all('./request/final/sortedd01.csv')

if __name__ == "__main__":
   main()

