import adafruit_dht
import time
import board
import base64
import requests
import datetime
import json
import uuid
from utils.jwt import create_jwt
from gcp_iot.request_builder import IOTInterface

# --------- User Settings ---------
SECONDS_BETWEEN_READS = 60
# ---------------------------------

def get_single_reading(dhtSensor):
    try:
        humidity = dhtSensor.humidity
        temp = dhtSensor.temperature
        humidity = format(humidity,".2f")
        temp = format(temp, ".2f")
    except RuntimeError as e:
        print(e)
        temp = "-999"
        humidity = "-999"
    #print("Temperature(C)", temp)
    #print("Humidity(%)", humidity)

    return temp, humidity


def structure_reading(temp, humidity) -> str:
    current_time = datetime.datetime.now()
    json_obj = {
        'temperature': temp,
        'humidity': humidity,
        'measurement_time': str(current_time),
        'device_message_id': str(uuid.uuid4())
    }
    return json.dumps(json_obj)


def measure_temperature():
    # Initialise sensor interface and GCP IOT interface
    dhtSensor = adafruit_dht.DHT22(board.D4)
    iot_interface = IOTInterface()

    while True:
        temp_c, humidity = get_single_reading(dhtSensor)

        structured_reading = structure_reading(temp_c, humidity)

        iot_interface.make_request(structured_reading)
        
        time.sleep(SECONDS_BETWEEN_READS)

if __name__ == "__main__":
    measure_temperature()