import requests
import json
import pandas as pd
import numpy as np
from sklearn.preprocessing import StandardScaler

data = ["0","0","2002","9","5","06:21","$134.09","Swipe Transaction","3527213246127876953","91750.0"]
df = pd.DataFrame(data, columns = ['Original Data'])
df['Preprocessed Data'] = ["0","0","2002","9","5","0621","134.09","Swipe Transaction","3527213246127876953","91750.0"]
# print(df)

data_to_preproces = {"data": df['Original Data'].values.tolist()}
# print(data_to_preproces)

response = requests.put('http://localhost:8080/ai_inference', json=data_to_preproces)

# df['API Preprocessed Data'] = json.loads(response.content.decode())
# print(df)

print('Input data:')
print(df['Original Data'].values.tolist())

print('\nOutput response:')
print(response.content.decode())

out_data = json.loads(response.content.decode())

if float(out_data['outputs'][0]['data'][0]) > 0:
    print('Prediction:\nFraud')
else:
    print('Prediction:\nNot Fraud')