import { Module } from '@nestjs/common';
import { FraudService } from './fraud/fraud.service';

@Module({
  providers: [FraudService],
})
export class InvoicesModule {}
