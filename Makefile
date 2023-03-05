deploybot:
	make -C svcbot deploy

deploymid:
	make -C svcmid deploy

deploytop:
	make -C svctop deploy

deployall: deploybot deploymid deploytop

undeploybot:
	-make -C svcbot undeploy

undeploymid:
	-make -C svcmid undeploy

undeploytop:
	-make -C svctop undeploy

undeployall: undeploybot undeploymid undeploytop

all: undeployall deployall