# AppStream2.0 Nvidia GPU Monitoring Toolkit
-------------
## 파일목록
- config.py: 설정파일, 향후 확장성 및 유지보수 고려
- make_report.py: prom_to_csv.py의 결과를 분석하기 위함
- make_target.py: prometheus의 file 기반 service discovery용(참조 : https://prometheus.io/docs/guides/file-sd/)
- prom_to_csv.py: prometheus에서 http restful api 사용하여 메트릭 데이터 조회 
- prometheus.yml: 설정파일
- windows_nvidia_exporter.go: Windows Nvidia-SMI exporter. build 후 AppStream Image Builder에서 서비스화(sc.exe)해야 함
