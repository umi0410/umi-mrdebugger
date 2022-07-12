FROM alpine
WORKDIR /mrdebugger
COPY mrdebugger ./mrdebugger
CMD ["./mrdebugger"]