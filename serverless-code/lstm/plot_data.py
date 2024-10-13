import torch
import torch.nn as nn
import seaborn as sns
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from readData import read_azure_2019

file_name = './dataset/azurefunctions-dataset2019/invocations_per_function_md.anon.d01.csv'
function_id = '5e98666f0af3fee4a98e670ab893ddf57046816b30775e3decd77a098d317e98'

in_data = read_azure_2019(file_name, function_id)
fig_size = plt.rcParams["figure.figsize"]
fig_size[0] = 15
fig_size[1] = 5
plt.rcParams["figure.figsize"] = fig_size

plt.title('Time vs Invocations')
plt.ylabel('Invocations')
plt.xlabel('Time (minute)')
plt.grid(True)
plt.autoscale(axis='x',tight=True)
plt.plot(in_data['Invocation'])
plt.show()


