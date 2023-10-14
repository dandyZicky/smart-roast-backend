for i in {1..10}  # Change the range (1..10) to the number of iterations you want
do
  echo "Iteration $i: Executing command..."
  mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/109 -m '101' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/108 -m '120' && mosquitto_pub -h broker.hivemq.com -t tes_deh/benar/110 -m '104'
done
 
