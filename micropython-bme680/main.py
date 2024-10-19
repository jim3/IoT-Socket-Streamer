from bme680 import *
from machine import SoftI2C, Pin
import time

# Initialize the I2C bus using SoftI2C
i2c = SoftI2C(scl=Pin(22), sda=Pin(21))

# Initialize the BME680 sensor
bme = BME680_I2C(i2c)

# Ensure the sensor is properly initialized
if not bme:
    print("Failed to initialize BME680 sensor")
    while True:
        pass

while True:
    # Read sensor values
    temperature = bme.temperature
    humidity = bme.humidity
    pressure = bme.pressure
    gas = bme.gas

    # Print sensor values 91 hPa, Gas: 312511 ohms
    print(f"Temperature: {temperature} °C, Humidity: {humidity} %, Pressure: {pressure} hPa, Gas: {gas} ohms")

    # Wait for 1 second before reading again
    time.sleep(1)
