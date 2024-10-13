##create a map

import plotly.express as px
import pandas as pd
import numpy as np


fields = ['SITE_ID','LATITUDE', 'LONGITUDE', 'TYPE']

df = pd.read_csv("site-optus-melbCBD-divided.csv", usecols=fields)

print(df)

df.dropna(
    axis=0,
    how='any',
    thresh=None,
    subset=None,
    inplace=True
)

color_scale = ['red']

size = []

size = [100 for i in range(144)] 
 
for i in range(len(size)):
    if i > 124:
        size[i] = 0

fig = px.scatter_mapbox(df, 
                        lat="LATITUDE", 
                        lon="LONGITUDE", 
                        color_continuous_scale='red',
                        size=size,
                        color = 'TYPE',
                        size_max=5,
                        zoom=10, 
                        height=500,
                        width=800,
                        color_discrete_map={'Edge Server': 'red'}
                        )

# the node to be connected in sequence
borders = [1,5,8,12,13,1,13,6,18,19,15,2,11,10,3,4,9]
borders = np.array(borders)
#convert it to dataframe index
borders = borders + 124
fig.add_traces(px.line_mapbox(df.loc[borders], lat="LATITUDE", lon="LONGITUDE", color = 'TYPE', color_discrete_map={'Border': 'black'}).data)
##fig.update_traces(showlegend=False)
fig.update_layout(mapbox_style="open-street-map")

##legend position
fig.update_layout(legend=dict(
    yanchor="top",
    y=0.90,
    xanchor="right",
    x=0.99
))

fig.update_layout(legend_title="", legend=dict(title_font_family="Arial",font=dict(size= 10), bgcolor = 'rgba(210,215,211,1)'))
##only show edge server in legend
labels_to_show_in_legend = ["Edge Server"]

for trace in fig['data']: 
    if (not trace['name'] in labels_to_show_in_legend):
        trace['showlegend'] = False


fig.show()