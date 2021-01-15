# quic-experiment

Um simples experimento usando QUIC, medindo a performance para transferência de dados em redes WiFi comparando com o TCP. 
Para a disciplina de Redes Móveis com o professor José Rezende. 


### Metodologia
Para realizar-se a análise comparativa, foram seguidos os seguintes passos:

1. Gerar arquivos de texto em memória de tamanhos pré-determinados (100KB, 1MB, 10MB, 100MB);
2. Calcular o hash MD5 do arquivo gerado, para que se possa conferir a integridade do arquivo posteriormente;
3. O lado do cliente inicia um contador de tempo e estabelece uma conexão QUIC ou TCP entre si e o servidor;
4. Empacotar o arquivo numa mensagem seguindo o protocolo de comunicação estabelecido. Onde a mensagem consiste de 8 bytes contendo o tamanho do arquivo a ser enviado e em seguida os bytes do arquivo de fato;
5. Servidor recebe e desempacota a mensagem, calculando então o hash MD5 da mesma. Para então responder o hash de volta para o cliente.
6. Ao receber a mensagem com o hash o cliente finaliza a contagem de tempo, armazenando para uma análise posterior dos dados. Então confere se o hash recebido do servidor é o mesmo hash calculado antes de enviar o arquivo;

Os arquivos gerados de maneira determinística, de modo que aquivos de mesmo tamanho são sempre os mesmos arquivos com os mesmos hashes.

Ambos os protocolos oferecem garantia de entrega de dados, o cálculo e resposta do hash é utilizado mais para confirmação de entrega e uma garantia extra da integridade da transmissão.

Cabe-se ressaltar que a proposta de transmissão é simples, não se utilizando de algumas características interessantes do QUIC, tais como o reestabelecimento de conexões anteriores em 0-RTT e o uso de multiplexação de streams para enviar múltiplos arquivos simultaneamente. De tal maneira que o uso do QUIC fica o mais parecido possível com o TCP.

Quanto a rede Wifi, foram feitos experimentos em 3 cenários com sinais de recepção diferentes, classificados como "Excelente", "Razoável" e "Fraco". É importante ressaltar que os experimentos foram feitos numa rede residencial, propensa a interferências externas.

### Resultados

Cenário com sinal excelente:

![Tabela Cenário Excelente](https://raw.githubusercontent.com/fogodev/quic-experiment/main/results/tabela_excelente.png)
![Gráfico Cenário Excelente](https://raw.githubusercontent.com/fogodev/quic-experiment/main/results/grafico_excelente.png)

Cenário com sinal razoável:

![Tabela Cenário Razoável](https://raw.githubusercontent.com/fogodev/quic-experiment/main/results/tabela_razoavel.png)
![Gráfico Cenário Razoável](https://raw.githubusercontent.com/fogodev/quic-experiment/main/results/grafico_razoavel.png)

Cenário com sinal fraco:

![Tabela Cenário Fraco](https://raw.githubusercontent.com/fogodev/quic-experiment/main/results/tabela_fraco.png)
![Gráfico Cenário Fraco](https://raw.githubusercontent.com/fogodev/quic-experiment/main/results/grafico_fraco.png)

Mesmo sem se utilizar nenhuma das características diferenciais do QUIC em relação ao TCP, é possível notar que a existência das mesmas não acarreta em nenhum custo extra de processamento. Tendo em vista que as médias de tempos de transmissão para cada arquivo sempre estão dentro dos desvios padrões para o outro protocolo nos cenários "Excelente" e "Razoável".

Dentre os três cenários propostos, o cenário com intensidade de sinal fraco é o mais suscetível a interferências externas, o que afeta a reprodutibilidade dos resultados.

