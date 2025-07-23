// Objeto que vai permitir transferencia de dados entre as camadas da aplicação
export class ProcessInvoiceFraudDto {
  invoice_id: string;
  account_id: string;
  amount: number;
}
