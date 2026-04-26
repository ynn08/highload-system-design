export enum DeliveryStatus {
  PENDING = 'PENDING',
  ASSIGNED = 'ASSIGNED',
  PICKED_UP = 'PICKED_UP',
  DELIVERING = 'DELIVERING',
  COMPLETED = 'COMPLETED',
}

export interface Location {
  lat: number;
  lon: number;
}

export interface Courier {
  id: string;
  name: string;
  location: Location;
  isAvailable: boolean;
  currentOrders: string[]; // List of order IDs currently being delivered
}

export interface Delivery {
  orderId: string;
  restaurantId: string;
  status: DeliveryStatus;
  courierId?: string;
  destination: Location;
}

export interface CourierRepository {
  findNearestAvailable(location: Location): Promise<Courier | null>;
  findNearestWithCapacity(location: Location, restaurantId: string): Promise<Courier | null>;
  updateLocation(courierId: string, location: Location): Promise<void>;
  assignOrder(courierId: string, orderId: string): Promise<void>;
}

export interface DeliveryRepository {
  save(delivery: Delivery): Promise<void>;
  findById(orderId: string): Promise<Delivery | null>;
}
