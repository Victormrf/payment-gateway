import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { InvoicesModule } from './invoices/invoices.module';
import { PrismaModule } from './prisma/prisma.module';
import { FraudService } from './invoices/fraud/fraud.service';

@Module({
  imports: [InvoicesModule, PrismaModule],
  controllers: [AppController],
  providers: [AppService, FraudService],
})
export class AppModule {}
