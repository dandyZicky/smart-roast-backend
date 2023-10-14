for i in {1..10}  # Change the range (1..10) to the number of iterations you want
do
  echo "Iteration $i: Executing command..."
  mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/1 -m '101' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/2 -m '120' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/3 -m '104' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/4 -m '101' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/5 -m '101' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/6 -m '101' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/7 -m '101' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/8 -m '101' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/9 -m '101' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/10 -m '101'
done

for i in {1..10}
do
  echo "deleting connection: $i"
  mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/1 -m '-1' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/2 -m '-1' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/3 -m '-1' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/4 -m '-1' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/5 -m '-1' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/6 -m '-1' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/7 -m '-1' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/8 -m '-1' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/9 -m '-1' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/10 -m '-1'
done
 
