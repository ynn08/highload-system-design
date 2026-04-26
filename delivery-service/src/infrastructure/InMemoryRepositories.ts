import { Courier, CourierRepository, Delivery, DeliveryRepository, Location } from "../domain/entities";

export class InMemoryCourierRepository implements CourierRepository {
  private couriers: Courier[] = [
    { id: '1', name: 'Ivan', location: { lat: 55.75, lon: 37.62 }, isAvailable: true, currentOrders: [] },
    { id: '2', name: 'Petr', location: { lat: 55.80, lon: 37.70 }, isAvailable: true, currentOrders: [] },
  ];

  async findNearestAvailable(location: Location): Promise<Courier | null> {
    let nearest: Courier | null = null;
    let minDistance = Infinity;

    for (const courier of this.couriers) {
      if (!courier.isAvailable || courier.currentOrders.length >= 3) continue;
      const dist = this.calculateDistance(courier.location, location);
      if (dist < minDistance) {
        minDistance = dist;
        nearest = courier;
      }
    }
    return nearest;
  }

  async findNearestWithCapacity(location: Location, restaurantId: string): Promise<Courier | null> {
    // In a real system, we'd check if any courier is already at this restaurantId
    // For this mock, let's say if a courier has orders, they are "at a restaurant"
    for (const courier of this.couriers) {
      if (courier.currentOrders.length > 0 && courier.currentOrders.length < 3) {
        // Simple logic: if they have capacity and are already working, they can batch
        return courier;
      }
    }
    return null;
  }

  async updateLocation(courierId: string, location: Location): Promise<void> {
    const courier = this.couriers.find(c => c.id === courierId);
    if (courier) courier.location = location;
  }

  async assignOrder(courierId: string, orderId: string): Promise<void> {
    const courier = this.couriers.find(c => c.id === courierId);
    if (courier) {
      courier.currentOrders.push(orderId);
    }
  }

  private calculateDistance(l1: Location, l2: Location): number {
    return Math.sqrt(Math.pow(l1.lat - l2.lat, 2) + Math.pow(l1.lon - l2.lon, 2));
  }
}

export class InMemoryDeliveryRepository implements DeliveryRepository {
  private deliveries = new Map<string, Delivery>();

  async save(delivery: Delivery): Promise<void> {
    this.deliveries.set(delivery.orderId, delivery);
  }

  async findById(orderId: string): Promise<Delivery | null> {
    return this.deliveries.get(orderId) || null;
  }
}
