# AppStream2.0 Nvidia GPU Monitoring Toolkit
-------------
## 파일목록
- config.py: 설정파일, 향후 확장성 및 유지보수 고려
- make_report.py: prom_to_csv.py의 결과를 분석하기 위함
- make_target.py: prometheus의 file 기반 service discovery용(참조 : https://prometheus.io/docs/guides/file-sd/)
- prom_to_csv.py: prometheus에서 http restful api 사용하여 메트릭 데이터 조회 
- prometheus.yml: 설정파일
- windows_nvidia_exporter.go: Windows Nvidia-SMI exporter. build 후 AppStream Image Builder에서 서비스화(sc.exe)해야 함

## Docker 실행 명령어
docker run     -p 9090:9090     -v /root/prometheus.yml:/etc/prometheus/prometheus.yml -v /root/targets.js
on:/etc/prometheus/target.json    prom/prometheus

## 설정방법
1. windows_nvidia_exporter.go를 Build하여 exe 파일로 변환 (go build -v windows_nvidia_exporter.go)
2. 해당 실행파일을 원하는 경로에 이동
3. sc.exe를 사용하여 서비스 생성: sc.exe create SERVICE win_nvidia_exporter binpath="[PATH]"
4. Service 정상 작동 확인 (ctrl+r, services.msc)
5. Image Assistant 사용하여 Image 화
6. AWS Credential 정보 설정
7. 모니터링 서버에서 Prometheus 기동  
```docker run     -p 9090:9090     -v /root/prometheus.yml:/etc/prometheus/prometheus.yml -v /root/targets.json:/etc/prometheus/target.json    prom/prometheus```
8. prom_to_csv.py (파일로 redirect 필요), make_report.py 수행

## 비고
해당 테스트는 Appstream2.0 G4dn, Windows 2016을 대상으로 테스트됨
