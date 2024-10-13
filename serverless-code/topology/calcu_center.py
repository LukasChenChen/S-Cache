##divide the area into 10 areas
#area 1: node site_id, 1,2,14,13
#area 2: 13,14,11,12
#area 3: 2,3,16,14
#area 4: 14,16,15,17
#area 5: 15, 17, 10, 11
#area 6: 3,4,16,18
#area 7: 16,18,10,9
#area 8: 4, 5, 6, 18
#area 9: 6,7,18,19
#area 10:19,7,8,9

import plotly.express as px
import pandas as pd

fields = ['SITE_ID','LATITUDE', 'LONGITUDE', 'TYPE']

df = pd.read_csv("site-optus-melbCBD-divided.csv", usecols=fields)

# only keep the border nodes
df.drop(df.index[0:125], axis=0, inplace=True)

print(df)


#node 1
lat = []
long = []
for index, row in df.iterrows():
    if row['SITE_ID'] in (1,2,14,13):
        lat.append(row['LATITUDE'])
        long.append(row['LONGITUDE'])

latitude = []
longitude = []
latitude.append(sum(lat)/len(lat))
longitude.append(sum(long)/len(long))

#node 2
lat = []
long = []
for index, row in df.iterrows():
    if row['SITE_ID'] in (13,14,11,12):
        lat.append(row['LATITUDE'])
        long.append(row['LONGITUDE'])

latitude.append(sum(lat)/len(lat))
longitude.append(sum(long)/len(long))

#node 3
lat = []
long = []
for index, row in df.iterrows():
    if row['SITE_ID'] in (2,3,16,14):
        lat.append(row['LATITUDE'])
        long.append(row['LONGITUDE'])

latitude.append(sum(lat)/len(lat))
longitude.append(sum(long)/len(long))

#node 4
lat = []
long = []
for index, row in df.iterrows():
    if row['SITE_ID'] in (14,16,15,17):
        lat.append(row['LATITUDE'])
        long.append(row['LONGITUDE'])

latitude.append(sum(lat)/len(lat))
longitude.append(sum(long)/len(long))


# node 5 
lat = []
long = []
for index, row in df.iterrows():
    if row['SITE_ID'] in (15, 17, 10, 11):
        lat.append(row['LATITUDE'])
        long.append(row['LONGITUDE'])

latitude.append(sum(lat)/len(lat))
longitude.append(sum(long)/len(long))

# node 6 
lat = []
long = []
for index, row in df.iterrows():
    if row['SITE_ID'] in (3,4,16,18):
        lat.append(row['LATITUDE'])
        long.append(row['LONGITUDE'])

latitude.append(sum(lat)/len(lat))
longitude.append(sum(long)/len(long))

# node 7
lat = []
long = []
for index, row in df.iterrows():
    if row['SITE_ID'] in (16,18,10,9):
        lat.append(row['LATITUDE'])
        long.append(row['LONGITUDE'])

latitude.append(sum(lat)/len(lat))
longitude.append(sum(long)/len(long))

# node 8
lat = []
long = []
for index, row in df.iterrows():
    if row['SITE_ID'] in (4, 5, 6, 18):
        lat.append(row['LATITUDE'])
        long.append(row['LONGITUDE'])

latitude.append(sum(lat)/len(lat))
longitude.append(sum(long)/len(long))

# node 9
lat = []
long = []
for index, row in df.iterrows():
    if row['SITE_ID'] in (6,7,18,19):
        lat.append(row['LATITUDE'])
        long.append(row['LONGITUDE'])

latitude.append(sum(lat)/len(lat))
longitude.append(sum(long)/len(long))

# node 10
lat = []
long = []
for index, row in df.iterrows():
    if row['SITE_ID'] in (19,7,8,9):
        lat.append(row['LATITUDE'])
        long.append(row['LONGITUDE'])

latitude.append(sum(lat)/len(lat))
longitude.append(sum(long)/len(long))
    
data = {'LATITUDE': latitude, 'LONGITUDE': longitude}
center_nodes = pd.DataFrame(data)
print(center_nodes)



size = [100 for i in range(10)] 
fig = px.scatter_mapbox(center_nodes, 
                        lat="LATITUDE", 
                        lon="LONGITUDE", 
                        color_continuous_scale='red',
                        size=size,
                        # color = 'LATITUDE',
                        size_max=5,
                        zoom=10, 
                        height=500,
                        width=800,
                        )

fig.update_layout(mapbox_style="open-street-map")
# fig.update_layout(margin={"r":0,"t":0,"l":0,"b":0})
fig.show()

center_nodes.to_csv('./pi_nodes.csv', sep=',')