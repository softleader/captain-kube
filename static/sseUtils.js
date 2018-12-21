const DEFAULT_TEXT_DECODER = new TextDecoder("utf-8");
function readSSE(response, {onNext = (value) => {}, onComplete = () => {}, onError = (e) => {console.error(e)}}) {
  const {body, headers} = response;
  const reader = body.getReader();
  const stream = new ReadableStream({
    start(controller) {
      return pump();

      function pump(storedValueTx = '') {
        return reader.read().then(({done, value}) => {
          let valueTxStore = storedValueTx;

          // When no more data needs to be consumed, close the stream
          if (done) {
            onComplete();
            controller.close();
            return;
          }
          try {
            if (value) {
              // 此處的Value受限於Uint8Array大小，每65536會截斷一次
              valueTxStore += DEFAULT_TEXT_DECODER.decode(value);

              let regexp = /data:(?<data>.*)\n(\n|data:)/g;
              let matcher = regexp.exec(valueTxStore);
              let lastIndex = 0;
              while (matcher) {
                let dataTx = matcher.groups['data'];
                // 紀錄最後一次讀到的地方
                lastIndex = matcher.index + matcher[0].length;
                if (dataTx) {
                  let data;
                  try {
                    data = JSON.parse(dataTx);
                  } catch (e) {
                    data = dataTx;
                  }

                  onNext(data);
                }
                matcher = regexp.exec(valueTxStore);
              }
              // 將最後一次讀到的段落之後擷出來，以便跟下一次組合在一起
              valueTxStore = valueTxStore.slice(lastIndex);
            }
          } catch (e) {
            onError(e);
            onComplete();
            controller.close();
            return;
          }
          // Enqueue the next data chunk into our target stream
          controller.enqueue(value);
          return pump(valueTxStore);
        });
      }
    }
  });

  return new Response(stream, {headers});
}
