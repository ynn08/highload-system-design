package com.example.orderservice.infrastructure.messaging;

import com.example.orderservice.domain.model.Order;
import com.example.orderservice.domain.repository.MessagePublisher;
import lombok.RequiredArgsConstructor;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class KafkaMessagePublisher implements MessagePublisher {
    private final KafkaTemplate<String, Object> kafkaTemplate;

    @Override
    public void publishOrderCreated(Order order) {
        kafkaTemplate.send("order-created", order);
    }
}
