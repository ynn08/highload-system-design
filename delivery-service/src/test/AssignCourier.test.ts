import { AssignCourierUseCase } from '../application/AssignCourierUseCase';
import { InMemoryCourierRepository, InMemoryDeliveryRepository } from '../infrastructure/InMemoryRepositories';

describe('AssignCourierUseCase', () => {
  it('should assign the nearest available courier', async () => {
    const courierRepo = new InMemoryCourierRepository();
    const deliveryRepo = new InMemoryDeliveryRepository();
    const useCase = new AssignCourierUseCase(courierRepo, deliveryRepo);

    // User at 55.75, 37.62 (Ivan is here)
    const result = await useCase.execute('order-123', { lat: 55.75, lon: 37.62 });

    expect(result).toBeDefined();
    expect(result?.courier.name).toBe('Ivan');
    expect(result?.delivery.orderId).toBe('order-123');
  });
});
