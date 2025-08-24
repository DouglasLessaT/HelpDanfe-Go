# Frontend HelpDanfe

Interface web moderna e responsiva para o sistema HelpDanfe, baseada no design do fsist.com.br.

## 🚀 Funcionalidades

- **Consulta de NFe**: Interface intuitiva para consultar NFes pela chave de acesso
- **Consulta de Boletos**: Consulta múltiplos boletos de uma vez
- **Download de Arquivos**: Baixar XML e gerar PDF da NFe
- **Histórico**: Visualizar consultas realizadas anteriormente
- **Configurações**: Personalizar URL da API e timeout
- **Responsivo**: Funciona perfeitamente em desktop, tablet e mobile

## 📁 Estrutura de Arquivos

```
web/
├── index.html          # Página principal
├── css/
│   └── style.css       # Estilos CSS
├── js/
│   └── app.js          # JavaScript da aplicação
└── README.md           # Esta documentação
```

## 🎨 Design

O frontend foi desenvolvido seguindo o design do fsist.com.br com:

- **Cores**: Paleta baseada em roxo (#9f6aca) e azul (#031b4e)
- **Tipografia**: Roboto (Google Fonts)
- **Ícones**: Font Awesome 6
- **Layout**: Menu lateral responsivo
- **Componentes**: Cards, tabelas, formulários modernos

## 🔧 Configuração

### Pré-requisitos

- Navegador moderno (Chrome, Firefox, Safari, Edge)
- API HelpDanfe rodando (padrão: http://localhost:8080)

### Instalação

1. Clone o repositório
2. Navegue até a pasta `web/`
3. Abra o arquivo `index.html` no navegador

### Configuração da API

1. Acesse a seção "Configurações" no menu lateral
2. Configure a URL da API (padrão: http://localhost:8080/api/v1)
3. Ajuste o timeout conforme necessário
4. Clique em "Salvar Configurações"

## 📱 Uso

### Consultar NFe

1. Acesse a seção "Consultar NFe"
2. Digite a chave de acesso (44 dígitos)
3. Opcionalmente, selecione um certificado digital
4. Clique em "Consultar NFe"
5. Visualize os dados da NFe e boletos vinculados
6. Use os botões para baixar XML ou gerar PDF

### Consultar Boletos

1. Acesse a seção "Consultar Boletos"
2. Digite os códigos dos boletos (um por linha)
3. Clique em "Consultar Boletos"
4. Visualize os resultados em tabela
5. Copie códigos de barras com um clique

### Histórico

1. Acesse a seção "Histórico"
2. Visualize todas as consultas realizadas
3. Use os filtros por data e tipo
4. Os dados são salvos localmente no navegador

## 🛠️ Tecnologias Utilizadas

- **HTML5**: Estrutura semântica
- **CSS3**: Estilos modernos com Flexbox e Grid
- **JavaScript ES6+**: Funcionalidades dinâmicas
- **Fetch API**: Comunicação com a API
- **LocalStorage**: Armazenamento local de histórico e configurações

## 📋 Recursos

### Validações

- Chave de acesso deve ter exatamente 44 dígitos
- Apenas números são aceitos na chave de acesso
- Validação de formulários antes do envio

### Feedback Visual

- Loading spinner durante requisições
- Modal de erro para problemas
- Badges coloridos para status
- Animações suaves

### Responsividade

- Menu lateral colapsável em mobile
- Layout adaptativo para diferentes telas
- Botões e formulários otimizados para touch

## 🔒 Segurança

- Validação client-side de dados
- Sanitização de inputs
- Tratamento de erros de API
- Timeout configurável para requisições

## 🚀 Melhorias Futuras

- [ ] Autenticação de usuários
- [ ] Tema escuro/claro
- [ ] Exportação de relatórios
- [ ] Notificações push
- [ ] Cache offline
- [ ] PWA (Progressive Web App)

## 🐛 Solução de Problemas

### API não responde

1. Verifique se a API está rodando
2. Confirme a URL nas configurações
3. Verifique o console do navegador para erros

### Certificado não funciona

1. Verifique se o arquivo é .p12 ou .pfx
2. Confirme a senha do certificado
3. Teste com certificado válido

### Histórico não aparece

1. Verifique se o localStorage está habilitado
2. Limpe o cache do navegador
3. Tente em modo incógnito

## 📞 Suporte

Para suporte técnico:
- Email: suporte@helpdanfe.com
- Documentação: https://docs.helpdanfe.com
- Issues: https://github.com/seu-usuario/helpdanfe-go/issues

## 📄 Licença

MIT License - veja o arquivo LICENSE para detalhes.
