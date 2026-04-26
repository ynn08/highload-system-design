package com.example.orderservice.application.usecase;

import com.example.orderservice.domain.model.Order;
import com.example.orderservice.domain.model.OrderStatus;
import com.example.orderservice.domain.repository.MessagePublisher;
import com.example.orderservice.domain.repository.OrderRepository;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.math.BigDecimal;
import java.util.Collections;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CreateOrderUseCaseTest {

    @Mock
    private OrderRepository orderRepository;

    @Mock
    private MessagePublisher messagePublisher;

    @InjectMocks
    private CreateOrderUseCase createOrderUseCase;

    @Test
    void shouldCreateOrderSuccessfully() {
        Order input = Order.builder()
                .customerId(UUID.randomUUID())
                .restaurantId(UUID.randomUUID())
                .totalAmount(new BigDecimal("100.00"))
                .items(Collections.singletonList(any())) // simplified
                .build();

        when(orderRepository.save(any(Order.class))).thenAnswer(i -> i.getArguments()[0]);

        Order result = createOrderUseCase.execute(input);

        assertThat(result.getOrderId()).isNotNull();
        assertThat(result.getStatus()).isEqualTo(OrderStatus.PENDING);
        assertThat(result.getEstimatedDelivery()).isAfter(result.getCreatedAt());

        verify(orderRepository, times(1)).save(any(Order.class));
        verify(messagePublisher, times(1)).publishOrderCreated(any(Order.class));
    }
}
