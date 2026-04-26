package com.example.orderservice.infrastructure.persistence;

import com.example.orderservice.domain.model.Order;
import com.example.orderservice.domain.model.OrderStatus;
import com.example.orderservice.domain.repository.OrderRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Component;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
interface JpaOrderRepository extends JpaRepository<OrderEntity, UUID> {}

@Component
@RequiredArgsConstructor
public class OrderRepositoryImpl implements OrderRepository {
    private final JpaOrderRepository jpaOrderRepository;

    @Override
    public Order save(Order order) {
        OrderEntity entity = new OrderEntity();
        entity.setOrderId(order.getOrderId());
        entity.setCustomerId(order.getCustomerId());
        entity.setRestaurantId(order.getRestaurantId());
        entity.setStatus(order.getStatus().name());
        entity.setTotalAmount(order.getTotalAmount());
        entity.setDeliveryFee(order.getDeliveryFee());
        entity.setSurgeMultiplier(order.getSurgeMultiplier());
        entity.setDeliveryAddress(order.getDeliveryAddress());
        entity.setEstimatedDelivery(order.getEstimatedDelivery());
        entity.setCreatedAt(order.getCreatedAt());
        
        jpaOrderRepository.save(entity);
        return order;
    }

    @Override
    public Optional<Order> findById(UUID id) {
        return jpaOrderRepository.findById(id).map(entity -> {
            Order order = new Order();
            order.setOrderId(entity.getOrderId());
            order.setCustomerId(entity.getCustomerId());
            order.setRestaurantId(entity.getRestaurantId());
            order.setStatus(OrderStatus.valueOf(entity.getStatus()));
            order.setTotalAmount(entity.getTotalAmount());
            order.setDeliveryFee(entity.getDeliveryFee());
            order.setSurgeMultiplier(entity.getSurgeMultiplier());
            order.setDeliveryAddress(entity.getDeliveryAddress());
            order.setEstimatedDelivery(entity.getEstimatedDelivery());
            order.setCreatedAt(entity.getCreatedAt());
            return order;
        });
    }
}
