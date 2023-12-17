#!/bin/bash
docker login -u hoanggiangco94 -p "glpat-y5xdxy7rNb5uCsEin_Af" "registry.gitlab.com/mmlabel/mm-printing-web"
docker build -t registry.gitlab.com/mmlabel/mm-printing-backend:${TAG} .
docker push registry.gitlab.com/mmlabel/mm-printing-backend:${TAG}

git tag ${TAG}
git push origin ${TAG}