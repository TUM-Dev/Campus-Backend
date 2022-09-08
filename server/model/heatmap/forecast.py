import sys
import sqlite3
import requests
import pandas as pd
from  prophet import Prophet
from datetime import date, datetime

def get_ap_names():
  conn = sqlite3.connect('./data/sqlite/heatmap.db')
  cur = conn.cursor()
  cur.execute("SELECT DISTINCT Name FROM apstat WHERE Address LIKE '%TUM%' AND Lat!='lat' AND Long!='long'")
  names = []
  for row in cur:
    names.append(row[0])
  return names


def get_json_from_AP(apName):
  eduroam     = "alias(ap.{ap}.ssid.eduroam,\"eduroam\")".format(ap = apName)
  lrz         = "alias(ap.{ap}.ssid.lrz,\"lrz\")".format(ap = apName)
  mwnEvents   = "alias(ap.{ap}.ssid.mwn-events,\"mwn-events\")".format(ap = apName)
  BayernWLAN  = "alias(ap.{ap}.ssid.@BayernWLAN,\"@BayernWLAN\")".format(ap = apName)
  other       = "alias(ap.{ap}.ssid.other,\"other\")".format(ap = apName)
  target      = "cactiStyle(group({eduroam},{lrz},{mwn_events},{bayern},{other}))".format(eduroam = eduroam, lrz = lrz, mwn_events = mwnEvents, bayern = BayernWLAN, other = other)
  url = "http://graphite-kom.srv.lrz.de/render?target={target}&format=json&from=-10days".format(target = target)
  jsonResp = requests.get(url).json()
  return jsonResp

# modify json: calculate total connections per timestamp
# change order timestamp <-> devices
# outputs data ready for forecasting
def process_json(jsonResp):
  total = []
  for datapoint in jsonResp[0]['datapoints']:
    connDev = datapoint[0]
    ts = datapoint[1]
    tm = datetime.utcfromtimestamp(ts).strftime('%Y-%m-%d %H:%M:%S')
    if connDev is None:
      total.append([ts, 0])
    else:
      total.append([ts, connDev])

  for jsonEntry in jsonResp[1:]:
    datapoints = jsonEntry['datapoints']
    for idx, dpArr in enumerate(datapoints):
      ts = dpArr[1]
      tm = datetime.utcfromtimestamp(ts).strftime('%Y-%m-%d %H:%M:%S')
      connDevices = dpArr[0]
      currDev = total[idx][1]
      if connDevices is not None:
        total[idx] = [tm, currDev + connDevices]
      else:
        total[idx] = [tm, currDev]
  return total

def forecast(total):
  df1 = pd.DataFrame(total)
  df1.columns = ['ds', 'y']
#   print(df1)

  df1.index.name = None
  m = Prophet()
  df1['floor'] = 0.0
  m.fit(df1)
  future = m.make_future_dataframe(periods=15*24, freq='H')
  future['floor'] = 0.0
  # future.tail()

  forecast = m.predict(future)
  
  dfout = forecast[['ds', 'trend']].tail(15)
  dfout['ds'] = pd.to_datetime(dfout['ds'])

  dfout.set_index('ds', inplace=True)
  
  return dfout #forecast.tail(15)

def process_forecasted_data(data):
    hourlyAvgs = []
    for row in data.iterrows():
        dateObj = pd.to_datetime(row[0])
        hour = dateObj.hour
        day = dateObj.day
        trend = row[1]['trend']
        hourlyAvgs.append((day, hour, trend))
    return hourlyAvgs


# processes data from csv file containing list of [timestamp, value]
# to a list of [day, hr, avg]
def process_forecasted_data_csv(forecastedData):
  col_list = ["ds", "trend"]
  df = pd.read_csv(forecastedData, usecols=col_list)

  hourlyAvgs = []
  prevHour = -1
  trendTot = 0
  cnt = 0
  for _, row in df.iterrows():
    dateObj = pd.to_datetime(row['ds'])
    hour = dateObj.hour
    day = dateObj.day
    trend = row['trend']
    if prevHour < 0:
      prevHour = hour
      trendTot = trend
      cnt = 1
      continue
    if hour != prevHour:
      # calculate trend avg. for prevHour
      hourlyAvgs.append((day, prevHour, trendTot / cnt))
      prevHour = hour
      trendTot = trend
      cnt = 1
    else:
      trendTot += trend
      cnt += 1
    minutes = dateObj.minute
    if hour == 23 and minutes == 45:
      hourlyAvgs.append((day, prevHour, trendTot / cnt))
      prevHour = -1

  return hourlyAvgs


# expects processed forecasted data in format [(day1, hr1, avg),(day1, hr2, avg)...]
# and stores it in the database 
def write_to_DB(forecastData, name):
  conn = sqlite3.connect('./data/sqlite/heatmap.db')
  cur = conn.cursor()
  for day, hour, avg in forecastData:
    stmt = ""
    if hour == 0:
      stmt = 'UPDATE future SET T0 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 1:
      stmt = 'UPDATE future SET T1 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 2:
      stmt = 'UPDATE future SET T2 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 3:
      stmt = 'UPDATE future SET T3 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 4:
      stmt = 'UPDATE future SET T4 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 5:
      stmt = 'UPDATE future SET T5 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 6:
      stmt = 'UPDATE future SET T6 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 7:
      stmt = 'UPDATE future SET T7 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 8:
      stmt = 'UPDATE future SET T8 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 9:
      stmt = 'UPDATE future SET T9 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 10:
      stmt = 'UPDATE future SET T10 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 11:
      stmt = 'UPDATE future SET T11 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 12:
      stmt = 'UPDATE future SET T12 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 13:
      stmt = 'UPDATE future SET T13 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 14:
      stmt = 'UPDATE future SET T14 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 15:
      stmt = 'UPDATE future SET T15 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 16:
      stmt = 'UPDATE future SET T16 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 17:
      stmt = 'UPDATE future SET T17 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 18:
      stmt = 'UPDATE future SET T18 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 19:
      stmt = 'UPDATE future SET T19 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 20:
      stmt = 'UPDATE future SET T20 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 21:
      stmt = 'UPDATE future SET T21 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 22:
      stmt = 'UPDATE future SET T22 = ? WHERE Day = ? AND AP_Name = ?'
    elif hour == 23:
      stmt = 'UPDATE future SET T23 = ? WHERE Day = ? AND AP_Name = ?'
    
    cur.execute(stmt, [avg, day, name])
  conn.commit()
  # conn.close()


def predictAndStoreInDB(names):
  for name in names:
    jsonData = get_json_from_AP(name)
    processedData = process_json(jsonData)
    forecastedData = forecast(processedData)
    processedForecast = process_forecasted_data(forecastedData)
    write_to_DB(processedForecast, name)


names = get_ap_names()
predictAndStoreInDB(names)