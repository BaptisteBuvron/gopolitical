# Gopolitical Server

## Percevoir

### Produire des ressources

Un channel est dédiée entre chaque pays et le serveur. Le pays envoie un message au serveur pour lui demander combien de ressource il à produit, combien de ressource il à vendu et combien de ressource il à acheté.

- Demander à l'environnement combien de ressource il à produire

### Récupérer des informations

- Si amis, informations à -10% +10%
- Si neutre, informations à -25% +25%
- Si enemies, informations à -50% +50%

## Deliberate

### Choix des echanges

Stock de sécurité = consomation * 10
Stock = ressourc actuelles

-5% si produit non-vendu
+5% si produit vendu

Un pays attaque si :
    - c'est moins chère pour lui d'attaquer

prix_achat = prix_du_marché(Stock de sécurité - Stock)

prix_voisin = prix_du_marché(Stock de sécurité du voisin) + production * CONSTANTE

prix_achat < prix_voisin





prix_du_marché(production)
prix d'achat actuelle du marché pour ces besoins < production du pays voisin * 

production
Somme()

### Choix des pays à attaquer

## Agir

### Questions

#### Gestion du marché

Un pays fait une demande d'achat ou de vente au marché, il indique également le territoire concerné par la demande. Le marché met à jour les ordres d'achat et de vente du territoire concerné.

- Lorsque l'environnement recoit une demande de ventes par un pays, il met à jour l'ordre d'achat du marché.
- Lorsque l'environnement recoit une demande d'achat par un pays, il met à jour l'ordre de vente du marché. 
  
A chaque fin de tour, l'environnement essaye de faire correspondre les ordres d'achat et de vente. Si un ordre d'achat est supérieur ou égal à un ordre de vente, alors l'environnement effectue la transaction et met à jour les stocks des pays concernés. Si un ordre d'achat est inférieur à un ordre de vente, alors l'environnement ne fait rien.

L'environnement stock ensuite les transactions effectuées dans un talbeau de transactions. Ce tableau est ensuite envoyé à chaque pays lors du tour suivant (la perception de chaque pays).