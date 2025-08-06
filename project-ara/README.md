# Project Ara - WhatsApp Financial Bot for Brazilian MEIs

A conversational WhatsApp bot that helps Brazilian MEIs manage their finances through text, voice, and image inputs.

## Project Overview

Project Ara is a financial management tool designed specifically for Brazilian Microempreendedores Individuais (MEIs). The bot allows users to:

- Log transactions via text, voice, or receipt photos
- Get real-time financial summaries
- Track income and expenses automatically
- Access a 50-transaction free trial before subscribing

## Features

### MVP Features (Phase 1-3)
- ✅ **Text Input Processing**: Natural language transaction logging
- ✅ **Voice Transcription**: Audio message processing
- ✅ **Image OCR**: Receipt photo analysis
- ✅ **Real-time Reporting**: Instant financial summaries
- ✅ **Trial Management**: 50-transaction free trial
- ✅ **Subscription System**: Paid plan conversion

### Technical Stack
- **Backend**: Go with Gin framework
- **Database**: PostgreSQL with GORM
- **Messaging**: WhatsApp Business API
- **AI Services**: OpenAI/Gemini (Phase 2)
- **Infrastructure**: Docker + AWS/Azure
- **Payments**: Brazilian gateways (Phase 3)

## Quick Start

### Prerequisites
- Go 1.21+
- Docker and Docker Compose
- PostgreSQL (for local development)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd project-ara
   ```

2. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

3. **Start with Docker Compose**
   ```bash
   docker-compose up --build
   ```

4. **Or run locally**
   ```bash
   # Start PostgreSQL
   docker run -d --name postgres \
     -e POSTGRES_DB=project_ara \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=password \
     -p 5432:5432 postgres:15-alpine

   # Run the application
   go run cmd/server/main.go
   ```

5. **Test the health endpoint**
   ```bash
   curl http://localhost:8080/health
   ```

## API Endpoints

### Health Check
- `GET /health` - Application health status

### WhatsApp Webhook
- `POST /api/v1/webhook/whatsapp` - WhatsApp message processing

### Transactions
- `POST /api/v1/transactions` - Create transaction
- `GET /api/v1/users/{id}/summary` - Get financial summary

### Subscriptions
- `POST /api/v1/subscriptions` - Create subscription

## Development Phases

### Phase 1: Foundation & Infrastructure ✅
- [x] Go project setup with Gin framework
- [x] PostgreSQL database with GORM
- [x] WhatsApp Business API integration
- [x] Docker containerization
- [x] Basic health check endpoints

### Phase 2: Core AI Services (Next)
- [ ] NLP service for text processing
- [ ] Voice transcription service
- [ ] OCR service for receipt analysis

### Phase 3: Business Logic & UX
- [ ] Transaction management
- [ ] Financial reporting
- [ ] Subscription & payment integration

### Phase 4: Integration & Testing
- [ ] End-to-end testing
- [ ] Load testing
- [ ] Security testing

### Phase 5: Deployment & Launch
- [ ] Production deployment
- [ ] Monitoring setup
- [ ] Launch preparation

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `PORT` | Server port | No (default: 8080) |
| `DB_HOST` | Database host | Yes |
| `DB_PORT` | Database port | No (default: 5432) |
| `DB_USER` | Database user | Yes |
| `DB_PASSWORD` | Database password | Yes |
| `DB_NAME` | Database name | Yes |
| `WHATSAPP_ACCESS_TOKEN` | WhatsApp API token | Yes |
| `WHATSAPP_PHONE_NUMBER_ID` | WhatsApp phone number ID | Yes |

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    trial_transactions_count INTEGER DEFAULT 0,
    subscription_status VARCHAR(20) DEFAULT 'trial'
);
```

### Transactions Table
```sql
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    amount DECIMAL(10,2) NOT NULL,
    description TEXT,
    transaction_type VARCHAR(10) NOT NULL,
    source VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    corrected_at TIMESTAMP,
    correction_data JSONB
);
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions, please contact the development team or create an issue in the repository. 