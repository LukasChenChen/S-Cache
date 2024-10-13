#using lstm to predict data from azure_2019 trace
#author: chen
#01.07.2022
import torch
import torch.nn as nn
import seaborn as sns
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from readData import read_azure_2019
#%matplotlib inline

#file name and function id
file_name = './dataset/azurefunctions-dataset2019/invocations_per_function_md.anon.d01.csv'
function_id = '5e98666f0af3fee4a98e670ab893ddf57046816b30775e3decd77a098d317e98'

in_data = read_azure_2019(file_name, function_id)

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

#normalize data to [-1,1]
from sklearn.preprocessing import MinMaxScaler

scaler = MinMaxScaler(feature_range=(-1, 1))
train_data_normalized = scaler.fit_transform(train_data .reshape(-1, 1))

train_data_normalized = torch.FloatTensor(train_data_normalized).view(-1)
#set the window to be 10minute
train_window = 10

## every 10000seconds
def create_inout_sequences(input_data, tw):
    inout_seq = []
    L = len(input_data)
    for i in range(L-tw):
        train_seq = input_data[i:i+tw]
        train_label = input_data[i+tw:i+tw+1]
        inout_seq.append((train_seq ,train_label))
    return inout_seq

train_inout_seq = create_inout_sequences(train_data_normalized, train_window)

train_inout_seq[:5]

class LSTM(nn.Module):
    def __init__(self, input_size=1, hidden_layer_size=100, output_size=1):
        super().__init__()
        self.hidden_layer_size = hidden_layer_size

        self.lstm = nn.LSTM(input_size, hidden_layer_size)

        self.linear = nn.Linear(hidden_layer_size, output_size)

        self.hidden_cell = (torch.zeros(1,1,self.hidden_layer_size),
                            torch.zeros(1,1,self.hidden_layer_size))

    def forward(self, input_seq):
        lstm_out, self.hidden_cell = self.lstm(input_seq.view(len(input_seq) ,1, -1), self.hidden_cell)
        predictions = self.linear(lstm_out.view(len(input_seq), -1))
        return predictions[-1]

model = LSTM()
loss_function = nn.MSELoss()
optimizer = torch.optim.Adam(model.parameters(), lr=0.001)

print(model)

epochs = 5

for i in range(epochs):
    for seq, labels in train_inout_seq:
        optimizer.zero_grad()
        model.hidden_cell = (torch.zeros(1, 1, model.hidden_layer_size),
                        torch.zeros(1, 1, model.hidden_layer_size))

        y_pred = model(seq)

        single_loss = loss_function(y_pred, labels)
        single_loss.backward()
        optimizer.step()

#    if i%25 == 1:
       # print(f'epoch: {i:3} loss: {single_loss.item():10.8f}')

#print(f'epoch: {i:3} loss: {single_loss.item():10.10f}')

fut_pred = 10

#test_inputs = train_data_normalized[-train_window:].tolist()
test_inputs = train_data_normalized[-test_data_size:].tolist()
print('test_inputs')
print(test_inputs)

model.eval()

for i in range(fut_pred):
    seq = torch.FloatTensor(test_inputs[-train_window:])
    with torch.no_grad():
        model.hidden = (torch.zeros(1, 1, model.hidden_layer_size),
                        torch.zeros(1, 1, model.hidden_layer_size))
        test_inputs.append(model(seq).item())

test_inputs[fut_pred:]

actual_predictions = scaler.inverse_transform(np.array(test_inputs[train_window:] ).reshape(-1, 1))
print(actual_predictions)

print('test dataset')
x = np.arange(1152, 1440, 1)
print(x)

plt.title('Time  vs Inovcation')
plt.ylabel('Invocations')
plt.grid(True)
plt.autoscale(axis='x', tight=True)
plt.plot(in_data['Invocation'], 'r')
plt.plot(x,actual_predictions, 'g')
#plt.plot(actual_predictions, 'g')
plt.show()




