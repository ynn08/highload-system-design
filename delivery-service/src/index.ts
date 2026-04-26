import express from 'express';
import { WebSocketServer } from 'ws';
import { Kafka } from 'kafkajs';
import http from 'http';
import { AssignCourierUseCase } from './application/AssignCourierUseCase';
import { InMemoryCourierRepository, InMemoryDeliveryRepository } from './infrastructure/InMemoryRepositories';

const app = express();
const server = http.createServer(app);
const wss = new WebSocketServer({ server });

const courierRepo = new InMemoryCourierRepository();
const deliveryRepo = new InMemoryDeliveryRepository();
const assignCourierUseCase = new AssignCourierUseCase(courierRepo, deliveryRepo);

const kafka = new Kafka({
  clientId: 'delivery-service',
  brokers: [process.env.KAFKA_BROKERS || 'localhost:9092']
});

const consumer = kafka.consumer({ groupId: 'delivery-group' });

const run = async () => {
  await consumer.connect();
  await consumer.subscribe({ topic: 'payment-processed', fromBeginning: true });

  await consumer.run({
    eachMessage: async ({ message }) => {
      if (!message.value) return;
      const event = JSON.parse(message.value.toString());
      
      if (event.status === 'PAID') {
        // Destination is normally in the event, using dummy for demo
        await assignCourierUseCase.execute(event.orderId, event.restaurantId, { lat: 55.75, lon: 37.62 });
        
        wss.clients.forEach((client) => {
          client.send(JSON.stringify({
            orderId: event.orderId,
            status: 'ASSIGNED',
          }));
        });
      }
    },
  });
};

run().catch(console.error);

app.get('/health', (req, res) => res.send('OK'));

app.get('/api/v1/orders/:orderId/status', async (req, res) => {
  const delivery = await deliveryRepo.findById(req.params.orderId);
  if (!delivery) {
    return res.status(404).json({ error: 'order not found' });
  }
  res.json(delivery);
});

const PORT = process.env.PORT || 8080;
server.listen(PORT, () => {
  console.log(`Delivery service listening on port ${PORT}`);
});
