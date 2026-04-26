package com.example.orderservice.infrastructure.messaging;

import com.example.orderservice.domain.model.OutboxEvent;
import com.example.orderservice.domain.repository.OutboxRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
@RequiredArgsConstructor
@Slf4j
public class OutboxProcessor {
    private final OutboxRepository outboxRepository;
    private final KafkaTemplate<String, Object> kafkaTemplate;

    @Scheduled(fixedDelay = 1000) // Every second
    public void processOutbox() {
        List<OutboxEvent> pendingEvents = outboxRepository.findPending();
        
        for (OutboxEvent event : pendingEvents) {
            try {
                log.info("Processing outbox event: {} for aggregate: {}", event.getId(), event.getAggregateId());
                
                // Determine topic based on event type (simplified)
                String topic = "order-created";
                if ("ORDER_UPDATED".equals(event.getEventType())) {
                    topic = "order-updated";
                }
                
                // Publish to Kafka
                kafkaTemplate.send(topic, event.getPayload()).get(); // Wait for confirmation
                
                // Mark as processed
                outboxRepository.markAsProcessed(event);
                
                log.info("Successfully processed outbox event: {}", event.getId());
            } catch (Exception e) {
                log.error("Failed to process outbox event: {}", event.getId(), e);
                // In production, we'd implement a retry count and eventually mark as FAILED
            }
        }
    }
}
