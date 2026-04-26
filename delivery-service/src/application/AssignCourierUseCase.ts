import { CourierRepository, DeliveryRepository, DeliveryStatus, Location } from '../domain/entities';

export class AssignCourierUseCase {
  constructor(
    private courierRepo: CourierRepository,
    private deliveryRepo: DeliveryRepository
  ) {}

  async execute(orderId: string, restaurantId: string, destination: Location) {
    console.log(`Assigning courier for order ${orderId} from restaurant ${restaurantId}...`);
    
    // 1. Try Batching: Find courier already at this restaurant with capacity
    let courier = await this.courierRepo.findNearestWithCapacity(destination, restaurantId);
    
    if (courier) {
      console.log(`Batching found! Courier ${courier.name} will pick up order ${orderId} alongside existing orders.`);
    } else {
      // 2. Fallback to nearest available
      courier = await this.courierRepo.findNearestAvailable(destination);
    }
    
    if (!courier) {
      console.warn(`No couriers available for order ${orderId}`);
      return;
    }

    const delivery = {
      orderId,
      restaurantId,
      status: DeliveryStatus.ASSIGNED,
      courierId: courier.id,
      destination
    };

    await this.deliveryRepo.save(delivery);
    await this.courierRepo.assignOrder(courier.id, orderId);
    
    console.log(`Courier ${courier.name} assigned to order ${orderId}`);
    
    return { delivery, courier };
  }
}
