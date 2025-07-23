import { repl } from '@nestjs/core';
import { AppModule } from './app.module';

async function bootstrap() {
  await repl(AppModule);
}
bootstrap();

// repl (read-eval-print loop) é um conceito muito comum em linguagens de programação interativas, como Python ou JavaScript, que permite aos desenvolvedores executar código em tempo real, testar expressões e explorar o ambiente de execução. No contexto do NestJS, o comando `repl` inicia um ambiente interativo onde você pode interagir com os módulos e serviços da sua aplicação NestJS.
