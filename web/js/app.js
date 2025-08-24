// Configurações globais
let config = {
    apiUrl: 'http://localhost:8080/api/v1',
    timeout: 30000
};

// Estado da aplicação
let currentNFe = null;
let currentBoletos = [];
let certificadosDisponiveis = [];
let certificadoSelecionado = null;

// Inicialização
document.addEventListener('DOMContentLoaded', function() {
    carregarConfiguracoes();
    carregarHistorico();
    setupEventListeners();
    verificarCertificadosDisponiveis();
});

// Configuração de event listeners
function setupEventListeners() {
    // Validação da chave de acesso
    document.getElementById('chave-acesso').addEventListener('input', function(e) {
        let value = e.target.value.replace(/\D/g, '');
        e.target.value = value;
    });

    // Auto-resize do textarea
    document.getElementById('codigos-boletos').addEventListener('input', function(e) {
        e.target.style.height = 'auto';
        e.target.style.height = e.target.scrollHeight + 'px';
    });
}

// Funções de certificado digital
async function verificarCertificadosDisponiveis() {
    try {
        // Verifica se o navegador suporta Web Crypto API
        if (!window.crypto || !window.crypto.subtle) {
            console.warn('Web Crypto API não suportada neste navegador');
            mostrarAviso('Certificados digitais não são suportados neste navegador');
            return;
        }

        // Atualiza a interface para mostrar status de certificados
        atualizarStatusCertificado('Verificando certificados disponíveis...');
        
        // Tenta acessar certificados via navigator.credentials (se disponível)
        if (navigator.credentials && navigator.credentials.get) {
            try {
                await verificarCredenciais();
            } catch (error) {
                console.log('Credenciais não disponíveis:', error.message);
            }
        }

        // Verifica se há certificados instalados no browser
        await verificarCertificadosNavigador();
        
    } catch (error) {
        console.error('Erro ao verificar certificados:', error);
        atualizarStatusCertificado('Erro ao verificar certificados');
    }
}

async function verificarCredenciais() {
    // Esta função tentará usar a API de credenciais se disponível
    // Nota: Esta API ainda é experimental em muitos navegadores
    try {
        const publicKey = {
            challenge: new Uint8Array(32),
            rp: { name: "HelpDanfe" },
            user: {
                id: new Uint8Array(16),
                name: "user@helpdanfe.com",
                displayName: "Usuário HelpDanfe"
            },
            pubKeyCredParams: [{ alg: -7, type: "public-key" }],
            authenticatorSelection: {
                authenticatorAttachment: "platform",
                userVerification: "required"
            }
        };
        
        // Esta é uma tentativa de usar credenciais, mas pode falhar
        // É mais um placeholder para futuras implementações
        console.log('Verificando credenciais WebAuthn...');
        
    } catch (error) {
        console.log('WebAuthn não disponível:', error.message);
    }
}

async function verificarCertificadosNavigador() {
    try {
        // Atualiza status
        atualizarStatusCertificado('Buscando certificados no sistema...');
        
        // Cria uma requisição de certificado cliente
        // Isso fará o navegador mostrar a lista de certificados disponíveis
        const response = await fetch(config.apiUrl + '/certificados/verificar', {
            method: 'GET',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        
        if (response.ok) {
            atualizarStatusCertificado('Certificados detectados automaticamente');
            document.getElementById('certificado-status').classList.add('certificado-disponivel');
        } else {
            atualizarStatusCertificado('Nenhum certificado detectado');
        }
        
    } catch (error) {
        console.log('Erro ao verificar certificados:', error);
        atualizarStatusCertificado('Clique para selecionar certificado');
    }
}

function atualizarStatusCertificado(mensagem) {
    const statusElement = document.getElementById('certificado-status');
    if (statusElement) {
        statusElement.textContent = mensagem;
    }
}

async function selecionarCertificado() {
    try {
        mostrarLoading();
        
        // Faz uma requisição que requer certificado cliente
        // Isso fará o navegador abrir o diálogo de seleção de certificado
        const response = await fetch(config.apiUrl + '/certificados/selecionar', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                acao: 'selecionar_certificado'
            })
        });
        
        if (response.ok) {
            const data = await response.json();
            certificadoSelecionado = data;
            atualizarStatusCertificado('Certificado selecionado: ' + (data.subject || 'Certificado válido'));
            document.getElementById('certificado-status').classList.add('certificado-selecionado');
            mostrarSucesso('Certificado selecionado com sucesso!');
        } else {
            mostrarErro('Erro ao selecionar certificado');
        }
        
    } catch (error) {
        console.error('Erro ao selecionar certificado:', error);
        mostrarErro('Erro ao acessar certificado. Verifique se há certificados instalados.');
    } finally {
        esconderLoading();
    }
}

async function usarCertificadoSistema() {
    try {
        // Esta função será chamada quando o usuário optar por usar certificado do sistema
        if (!certificadoSelecionado) {
            await selecionarCertificado();
        }
        
        return certificadoSelecionado;
    } catch (error) {
        console.error('Erro ao usar certificado do sistema:', error);
        throw error;
    }
}

// Funções do menu
function menu() {
    const menuLeft = document.getElementById('menu-left');
    const content = document.getElementById('content');
    
    if (menuLeft.className === 'menu-left-fechado') {
        menuLeft.className = 'menu-left-aberto';
        content.className = 'content-aberto';
    } else {
        menuLeft.className = 'menu-left-fechado';
        content.className = 'content-fechado';
    }
}

function showSection(sectionId) {
    // Esconder todas as seções
    document.querySelectorAll('.section').forEach(section => {
        section.classList.remove('active');
    });
    
    // Mostrar a seção selecionada
    document.getElementById(sectionId).classList.add('active');
    
    // Atualizar menu ativo
    document.querySelectorAll('.menu-item').forEach(item => {
        item.classList.remove('active');
    });
    
    // Encontrar e ativar o item do menu correspondente
    const menuItems = document.querySelectorAll('.menu-item');
    for (let item of menuItems) {
        if (item.getAttribute('onclick').includes(sectionId)) {
            item.classList.add('active');
            break;
        }
    }
}

// Funções de consulta de NFe
async function consultarNFe(event) {
    event.preventDefault();
    
    const chaveAcesso = document.getElementById('chave-acesso').value;
    const usarCertificadoArquivo = document.getElementById('usar-certificado-arquivo').checked;
    const certificado = document.getElementById('certificado').files[0];
    const senha = document.getElementById('senha-certificado').value;
    
    if (!chaveAcesso || chaveAcesso.length !== 44) {
        mostrarErro('A chave de acesso deve ter exatamente 44 dígitos');
        return;
    }
    
    mostrarLoading();
    
    try {
        let requestOptions;
        
        if (usarCertificadoArquivo && certificado) {
            // Usa certificado de arquivo (método tradicional)
            const formData = new FormData();
            formData.append('chave_acesso', chaveAcesso);
            formData.append('certificado', certificado);
            formData.append('senha', senha || '');
            formData.append('tipo_certificado', 'arquivo');
            
            requestOptions = {
                method: 'POST',
                body: formData
            };
        } else {
            // Usa certificado do sistema (novo método)
            try {
                await usarCertificadoSistema();
            } catch (certError) {
                mostrarErro('Erro ao acessar certificado do sistema: ' + certError.message);
                return;
            }
            
            requestOptions = {
                method: 'POST',
                credentials: 'include', // Importante: inclui certificados do navegador
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    chave_acesso: chaveAcesso,
                    tipo_certificado: 'sistema'
                })
            };
        }
        
        const response = await fetch(`${config.apiUrl}/nfe/consultar`, requestOptions);
        
        const data = await response.json();
        
        if (data.success) {
            currentNFe = data.data;
            exibirDadosNFe(data.data);
            await consultarBoletosNFe(chaveAcesso);
            adicionarAoHistorico('nfe', chaveAcesso, 'Consulta realizada com sucesso');
        } else {
            mostrarErro(data.error || 'Erro ao consultar NFe');
        }
    } catch (error) {
        console.error('Erro:', error);
        mostrarErro('Erro de conexão com a API');
    } finally {
        esconderLoading();
    }
}

function exibirDadosNFe(nfe) {
    const container = document.getElementById('dados-nfe');
    
    const html = `
        <div class="card">
            <div class="card-header">
                <h4 class="card-title">Informações Básicas</h4>
                <span class="badge badge-success">${nfe.status}</span>
            </div>
            <div class="row">
                <div class="col-md-6">
                    <p><strong>Chave de Acesso:</strong> ${nfe.chave_acesso}</p>
                    <p><strong>Número:</strong> ${nfe.numero}</p>
                    <p><strong>Série:</strong> ${nfe.serie}</p>
                    <p><strong>Data de Emissão:</strong> ${formatarData(nfe.data_emissao)}</p>
                </div>
                <div class="col-md-6">
                    <p><strong>Ambiente:</strong> ${nfe.ambiente}</p>
                    <p><strong>UF:</strong> ${nfe.uf}</p>
                    <p><strong>Valor Total:</strong> R$ ${formatarValor(nfe.valor_total)}</p>
                </div>
            </div>
        </div>
        
        <div class="card">
            <div class="card-header">
                <h4 class="card-title">Emitente</h4>
            </div>
            <p><strong>CNPJ:</strong> ${nfe.emitente_cnpj}</p>
            <p><strong>Nome:</strong> ${nfe.emitente_nome}</p>
            <p><strong>IE:</strong> ${nfe.emitente_ie}</p>
        </div>
        
        <div class="card">
            <div class="card-header">
                <h4 class="card-title">Destinatário</h4>
            </div>
            <p><strong>CNPJ:</strong> ${nfe.destinatario_cnpj}</p>
            <p><strong>Nome:</strong> ${nfe.destinatario_nome}</p>
            <p><strong>IE:</strong> ${nfe.destinatario_ie}</p>
        </div>
        
        <div class="card">
            <div class="card-header">
                <h4 class="card-title">Valores</h4>
            </div>
            <p><strong>Valor Total:</strong> R$ ${formatarValor(nfe.valor_total)}</p>
            <p><strong>Valor dos Produtos:</strong> R$ ${formatarValor(nfe.valor_produtos)}</p>
            <p><strong>Valor dos Impostos:</strong> R$ ${formatarValor(nfe.valor_impostos)}</p>
        </div>
    `;
    
    container.innerHTML = html;
    document.getElementById('resultado-nfe').style.display = 'block';
}

async function consultarBoletosNFe(chaveAcesso) {
    try {
        const response = await fetch(`${config.apiUrl}/nfe/${chaveAcesso}/boletos`);
        const data = await response.json();
        
        if (data.success) {
            currentBoletos = data.data.boletos || [];
            exibirBoletosNFe(currentBoletos);
        }
    } catch (error) {
        console.error('Erro ao consultar boletos:', error);
    }
}

function exibirBoletosNFe(boletos) {
    const container = document.getElementById('dados-boletos');
    
    if (boletos.length === 0) {
        container.innerHTML = '<p class="text-center">Nenhum boleto encontrado para esta NFe.</p>';
    } else {
        const html = `
            <table class="table">
                <thead>
                    <tr>
                        <th>Banco</th>
                        <th>Número</th>
                        <th>Valor</th>
                        <th>Vencimento</th>
                        <th>Status</th>
                        <th>Ações</th>
                    </tr>
                </thead>
                <tbody>
                    ${boletos.map(boleto => `
                        <tr>
                            <td>${boleto.banco}</td>
                            <td>${boleto.numero}</td>
                            <td>R$ ${formatarValor(boleto.valor)}</td>
                            <td>${formatarData(boleto.vencimento)}</td>
                            <td><span class="badge badge-${getStatusBadge(boleto.status)}">${boleto.status}</span></td>
                            <td>
                                <button class="btn btn-sm btn-secondary" onclick="copiarCodigoBarras('${boleto.codigo_barras}')">
                                    <i class="fas fa-copy"></i>
                                </button>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
        
        container.innerHTML = html;
    }
    
    document.getElementById('boletos-nfe').style.display = 'block';
}

// Funções de consulta de boletos
async function consultarBoletos(event) {
    event.preventDefault();
    
    const codigos = document.getElementById('codigos-boletos').value
        .split('\n')
        .map(codigo => codigo.trim())
        .filter(codigo => codigo.length > 0);
    
    if (codigos.length === 0) {
        mostrarErro('Digite pelo menos um código de boleto');
        return;
    }
    
    mostrarLoading();
    
    try {
        const response = await fetch(`${config.apiUrl}/boletos/consultar`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ codigos })
        });
        
        const data = await response.json();
        
        if (data.success) {
            exibirBoletosConsulta(data.data);
            adicionarAoHistorico('boleto', codigos.join(', '), `${data.data.length} boletos consultados`);
        } else {
            mostrarErro(data.error || 'Erro ao consultar boletos');
        }
    } catch (error) {
        console.error('Erro:', error);
        mostrarErro('Erro de conexão com a API');
    } finally {
        esconderLoading();
    }
}

function exibirBoletosConsulta(boletos) {
    const container = document.getElementById('dados-boletos-consulta');
    
    if (boletos.length === 0) {
        container.innerHTML = '<p class="text-center">Nenhum boleto encontrado.</p>';
    } else {
        const html = `
            <table class="table">
                <thead>
                    <tr>
                        <th>Banco</th>
                        <th>Número</th>
                        <th>Código de Barras</th>
                        <th>Valor</th>
                        <th>Vencimento</th>
                        <th>Status</th>
                    </tr>
                </thead>
                <tbody>
                    ${boletos.map(boleto => `
                        <tr>
                            <td>${boleto.banco}</td>
                            <td>${boleto.numero}</td>
                            <td>
                                <code>${boleto.codigo_barras}</code>
                                <button class="btn btn-sm btn-secondary" onclick="copiarCodigoBarras('${boleto.codigo_barras}')">
                                    <i class="fas fa-copy"></i>
                                </button>
                            </td>
                            <td>R$ ${formatarValor(boleto.valor)}</td>
                            <td>${formatarData(boleto.vencimento)}</td>
                            <td><span class="badge badge-${getStatusBadge(boleto.status)}">${boleto.status}</span></td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
        
        container.innerHTML = html;
    }
    
    document.getElementById('resultado-boletos').style.display = 'block';
}

// Funções de download
async function baixarXML() {
    if (!currentNFe) {
        mostrarErro('Nenhuma NFe consultada');
        return;
    }
    
    try {
        const response = await fetch(`${config.apiUrl}/nfe/${currentNFe.chave_acesso}/xml`);
        
        if (response.ok) {
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `nfe_${currentNFe.chave_acesso}.xml`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);
        } else {
            mostrarErro('Erro ao baixar XML');
        }
    } catch (error) {
        console.error('Erro:', error);
        mostrarErro('Erro ao baixar XML');
    }
}

async function gerarPDF() {
    if (!currentNFe) {
        mostrarErro('Nenhuma NFe consultada');
        return;
    }
    
    mostrarLoading();
    
    try {
        const response = await fetch(`${config.apiUrl}/nfe/${currentNFe.chave_acesso}/pdf`);
        
        if (response.ok) {
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `danfe_${currentNFe.chave_acesso}.pdf`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);
        } else {
            mostrarErro('Erro ao gerar PDF');
        }
    } catch (error) {
        console.error('Erro:', error);
        mostrarErro('Erro ao gerar PDF');
    } finally {
        esconderLoading();
    }
}

// Funções de utilidade
function copiarCodigoBarras(codigo) {
    navigator.clipboard.writeText(codigo).then(() => {
        mostrarSucesso('Código de barras copiado!');
    }).catch(() => {
        mostrarErro('Erro ao copiar código de barras');
    });
}

function formatarData(dataString) {
    if (!dataString) return '-';
    const data = new Date(dataString);
    return data.toLocaleDateString('pt-BR');
}

function formatarValor(valor) {
    if (!valor) return '0,00';
    return parseFloat(valor).toFixed(2).replace('.', ',');
}

function getStatusBadge(status) {
    switch (status.toLowerCase()) {
        case 'aberto':
        case 'pendente':
            return 'warning';
        case 'pago':
        case 'liquidado':
            return 'success';
        case 'vencido':
        case 'cancelado':
            return 'danger';
        default:
            return 'info';
    }
}

// Funções de formulário
function toggleCertificadoArquivo() {
    const checkbox = document.getElementById('usar-certificado-arquivo');
    const container = document.getElementById('certificado-arquivo-container');
    
    if (checkbox.checked) {
        container.style.display = 'block';
        container.classList.add('fade-in');
    } else {
        container.style.display = 'none';
        container.classList.remove('fade-in');
    }
}

function limparFormulario() {
    document.getElementById('nfe-form').reset();
    document.getElementById('resultado-nfe').style.display = 'none';
    document.getElementById('boletos-nfe').style.display = 'none';
    document.getElementById('certificado-arquivo-container').style.display = 'none';
    document.getElementById('usar-certificado-arquivo').checked = false;
    currentNFe = null;
    currentBoletos = [];
    // Reverifica certificados após limpar
    verificarCertificadosDisponiveis();
}

function limparFormularioBoletos() {
    document.getElementById('boletos-form').reset();
    document.getElementById('resultado-boletos').style.display = 'none';
}

// Funções de histórico
function adicionarAoHistorico(tipo, chave, resultado) {
    const historico = JSON.parse(localStorage.getItem('historico') || '[]');
    
    historico.unshift({
        id: Date.now(),
        tipo,
        chave,
        resultado,
        data: new Date().toISOString()
    });
    
    // Manter apenas os últimos 50 registros
    if (historico.length > 50) {
        historico.splice(50);
    }
    
    localStorage.setItem('historico', JSON.stringify(historico));
    carregarHistorico();
}

function carregarHistorico() {
    const historico = JSON.parse(localStorage.getItem('historico') || '[]');
    const container = document.getElementById('historico-lista');
    
    if (historico.length === 0) {
        container.innerHTML = '<p class="text-center">Nenhum histórico disponível.</p>';
        return;
    }
    
    const html = historico.map(item => `
        <div class="card">
            <div class="card-header">
                <h5 class="card-title">
                    <i class="fas fa-${item.tipo === 'nfe' ? 'file-invoice' : 'barcode'}"></i>
                    ${item.tipo.toUpperCase()}
                </h5>
                <small>${formatarData(item.data)}</small>
            </div>
            <p><strong>Chave/Código:</strong> ${item.chave}</p>
            <p><strong>Resultado:</strong> ${item.resultado}</p>
        </div>
    `).join('');
    
    container.innerHTML = html;
}

function filtrarHistorico() {
    const filtroData = document.getElementById('filtro-data').value;
    const filtroTipo = document.getElementById('filtro-tipo').value;
    
    let historico = JSON.parse(localStorage.getItem('historico') || '[]');
    
    if (filtroData) {
        historico = historico.filter(item => {
            const itemData = new Date(item.data).toDateString();
            const filtroDataObj = new Date(filtroData).toDateString();
            return itemData === filtroDataObj;
        });
    }
    
    if (filtroTipo) {
        historico = historico.filter(item => item.tipo === filtroTipo);
    }
    
    const container = document.getElementById('historico-lista');
    
    if (historico.length === 0) {
        container.innerHTML = '<p class="text-center">Nenhum resultado encontrado.</p>';
        return;
    }
    
    const html = historico.map(item => `
        <div class="card">
            <div class="card-header">
                <h5 class="card-title">
                    <i class="fas fa-${item.tipo === 'nfe' ? 'file-invoice' : 'barcode'}"></i>
                    ${item.tipo.toUpperCase()}
                </h5>
                <small>${formatarData(item.data)}</small>
            </div>
            <p><strong>Chave/Código:</strong> ${item.chave}</p>
            <p><strong>Resultado:</strong> ${item.resultado}</p>
        </div>
    `).join('');
    
    container.innerHTML = html;
}

// Funções de configuração
function carregarConfiguracoes() {
    const savedConfig = localStorage.getItem('config');
    if (savedConfig) {
        config = { ...config, ...JSON.parse(savedConfig) };
    }
    
    document.getElementById('api-url').value = config.apiUrl;
    document.getElementById('timeout').value = config.timeout / 1000;
}

function salvarConfiguracoes() {
    config.apiUrl = document.getElementById('api-url').value;
    config.timeout = parseInt(document.getElementById('timeout').value) * 1000;
    
    localStorage.setItem('config', JSON.stringify(config));
    mostrarSucesso('Configurações salvas com sucesso!');
}

function restaurarConfiguracoes() {
    config = {
        apiUrl: 'http://localhost:8080/api/v1',
        timeout: 30000
    };
    
    document.getElementById('api-url').value = config.apiUrl;
    document.getElementById('timeout').value = config.timeout / 1000;
    
    localStorage.removeItem('config');
    mostrarSucesso('Configurações restauradas!');
}

// Funções de UI
function mostrarLoading() {
    document.getElementById('loading').style.display = 'flex';
}

function esconderLoading() {
    document.getElementById('loading').style.display = 'none';
}

function mostrarErro(mensagem) {
    document.getElementById('erro-mensagem').textContent = mensagem;
    document.getElementById('modal-erro').style.display = 'flex';
}

function mostrarSucesso(mensagem) {
    // Implementar toast de sucesso
    alert(mensagem);
}

function mostrarAviso(mensagem) {
    console.warn(mensagem);
    // Implementar toast de aviso
    alert('Aviso: ' + mensagem);
}

function fecharModal() {
    document.getElementById('modal-erro').style.display = 'none';
}

// Fechar modal ao clicar fora
window.onclick = function(event) {
    const modal = document.getElementById('modal-erro');
    if (event.target === modal) {
        modal.style.display = 'none';
    }
}

// Fechar modal com ESC
document.addEventListener('keydown', function(event) {
    if (event.key === 'Escape') {
        document.getElementById('modal-erro').style.display = 'none';
    }
});
