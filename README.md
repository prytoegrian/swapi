# SwAPI

Cette application a été réalisée dans le contexte d'une évaluation des qualifications en golang.

## Évaluation et compilation
J'ai mis à disposition l'exécutable dans la release, mais une compilation locale est bien sûr possible. Après avoir importé le repo, placez-vous dans le répertoire racine de l'applicatif. Ensuite, exécutez `make install` pour installer le logiciel.

NB : Le `GOPATH` doit être configuré.

## Usage
Dans un terminal, lancez `swapi` pour faire tourner le serveur.
Dans un autre terminal, vous pourrez interroger le serveur aux deux routes disponibles :
* `GET, POST, OPTIONS` http://localhost:8080/peoples
* `GET, PUT, DELETE, OPTIONS` http://localhost:8080/peoples/{id:[0-9]+}


Le meilleur moyen pour le faire est de passer par `curl` :
```sh
curl -X GET http://localhost:8080/peoples
```

Comme attendu, cette route affiche la liste des personnages embarquant les véhicules et vaisseaux spatiaux du personnages. Il en sera de même pour la route `/peoples/ID`.

Les méthodes avec données `POST` et `PUT` doivent en plus définir une donnée via l'attribut `-d` :
```sh
curl -X POST -d '{"name": "Captain Planet", "height": 0, "mass": 0,  "hair": "unknown", "skin": "unknown", "eye": "unknown", "birth_year": "unknown", "gender": "female", "homeworld": 28, "films": "", "species": "", "vehicles": [], "starships": [], "url": "/captain"}' http://localhost:8080/peoples
```


## Choix techniques
Il n'était pas nécessaire de faire compliqué en terme de design. Le fichier `main.go` liste les routes possibles tandis que le fichier `handlers/handlers.go` les décrit, une à une. Les réponses obéissent au format [jsend](https://github.com/omniti-labs/jsend) afin de garantir une réponse normalisée aux clients.  
Bien que non requis, j'ai préféré séparer les différents aspects métiers en différents packages : `people`, `vehicle` et `starship`.

J'ai visé la testabilité, les packages métiers sont donc testés autant que possible. De la même façon, l'application est documentée selon les standards golang (au besoin, `make hard-lint` pour vérifier le respect des standards)
