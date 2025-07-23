import { Injectable } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { ProcessInvoiceFraudDto } from '../dto/process-invoice-fraud.dto';
import { Account, FraudReason, InvoiceStatus } from '@prisma/client';

@Injectable()
export class FraudService {
  constructor(private prismaService: PrismaService) {}

  async processInvoice(processInvoiceFraudDto: ProcessInvoiceFraudDto) {
    const { invoice_id, account_id, amount } = processInvoiceFraudDto;

    const invoice = await this.prismaService.invoice.findUnique({
      where: {
        id: invoice_id,
      },
    });

    if (invoice) {
      throw new Error('Invoice has already been processed');
    }

    // devemos agora ver se a conta ja existe ou deve ser criada
    // upsert - junção do insert e update
    const account = await this.prismaService.account.upsert({
      where: {
        id: account_id,
      },
      update: {},
      create: {
        id: account_id,
      },
    });

    const fraudResult = await this.detectFraud({ account, amount });

    await this.prismaService.invoice.create({
      data: {
        id: invoice_id,
        accountId: account.id,
        amount,
        ...(fraudResult && {
          fraudHistory: {
            create: {
              reason: fraudResult.reason!,
              description: fraudResult.description,
            },
          },
        }),
        status: fraudResult.hasFraud
          ? InvoiceStatus.REJECTED
          : InvoiceStatus.APPROVED,
      },
    });

    return {
      invoice,
      fraudResult,
    };
  }

  async detectFraud(data: { account: Account; amount: number }) {
    const { account, amount } = data;

    // check 1: verificar se a conta é suspeita
    if (account.isSuspicious) {
      return {
        hasFraud: true,
        reason: FraudReason.SUSPICIOUS_ACCOUNT,
        description: 'Account is suspicious',
      };
    }

    //
    const previuousInvoices = await this.prismaService.invoice.findMany({
      where: {
        accountId: account.id,
      },
      orderBy: {
        createdAt: 'desc',
      },
      take: 20, // pegar os ultimos 20 invoices
    });

    if (previuousInvoices.length) {
      const totalAmount = previuousInvoices.reduce((acc, invoice) => {
        return acc + invoice.amount;
      }, 0);

      const averageAmount = totalAmount / previuousInvoices.length;

      if (amount > averageAmount * (1 + 50 / 100) + averageAmount) {
        return {
          hasFraud: true,
          reason: FraudReason.UNUSUAL_PATTERN,
          description: `Amount ${amount} is higher than the avarege amount ${averageAmount} by more than 50%`,
        };
      }
    }

    const recentDate = new Date();
    recentDate.setDate(recentDate.getHours() - 24);

    const recentInvoices = await this.prismaService.invoice.findMany({
      where: {
        accountId: account.id,
        createdAt: {
          gte: recentDate,
        },
      },
    });

    if (recentInvoices.length > 100) {
      return {
        hasFraud: true,
        reason: FraudReason.FREQUENT_HIGH_VALUE,
        description: `Account ${account.id} has more than 100 invoices in the last 24 hours`,
      };
    }

    return {
      hasFraud: false,
    };
  }
}
