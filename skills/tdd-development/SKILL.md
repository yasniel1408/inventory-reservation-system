---
name: tdd-development
description: Usar cuando un developer implementa una feature, bugfix, refactor o cambio de comportamiento antes de escribir codigo productivo.
---

# TDD Development

Esta skill es obligatoria para agentes `developer` cuando implementan codigo productivo. La regla base es simple: primero se entiende el comportamiento, luego se escribe un test que falla, despues se implementa lo minimo para hacerlo pasar y recien ahi se refactoriza.

## Ciclo Obligatorio

```text
entender requerimiento
  -> elegir el comportamiento mas chico
  -> escribir test
  -> ejecutar y ver fallar por la razon esperada
  -> escribir codigo minimo
  -> ejecutar y ver pasar
  -> refactorizar sin cambiar comportamiento
  -> ejecutar tests otra vez
  -> repetir
```

## Reglas

- No escribir codigo productivo antes de tener un test fallando.
- El test debe probar comportamiento observable, no detalles internos.
- El test debe fallar por la razon esperada: feature ausente, bug reproducido o comportamiento incorrecto.
- Si el test pasa al primer intento, no sirve como RED; corregir el test o elegir otro comportamiento.
- En GREEN, escribir solo el codigo minimo necesario para pasar.
- Refactorizar solo despues de tener tests verdes.
- Si se descubre un bug, primero escribir un test que lo reproduzca.
- Si una parte no se puede automatizar todavia, pedir aprobacion explicita antes de saltar TDD y documentar el motivo.

## Evidencia Requerida

Todo `developer` debe reportar:

- Test escrito.
- Comando ejecutado para RED.
- Resultado/fallo esperado observado.
- Codigo minimo implementado.
- Comando ejecutado para GREEN.
- Refactor realizado, si aplica.
- Tests finales ejecutados.

## Excepciones

Solo se puede omitir TDD con aprobacion explicita del usuario o del `team-leader` cuando el trabajo sea:

- Documentacion pura.
- Configuracion sin comportamiento testeable.
- Scaffold mecanico sin logica.
- Spike descartable que no queda como codigo final.

Si el spike produce codigo final, ese codigo debe rehacerse con TDD.

## Checklist de Cierre

- [ ] Cada cambio de comportamiento tiene test.
- [ ] Se vio fallar el test antes del codigo productivo.
- [ ] El fallo RED fue por la razon esperada.
- [ ] El codigo GREEN fue minimo.
- [ ] Todos los tests relevantes pasan.
- [ ] No quedaron TODOs esenciales ni comportamiento sin cubrir.
