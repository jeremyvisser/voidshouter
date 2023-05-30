# Void Shouter

Are you frustrated by an insufficient number of messaging apps installed? Do you like unreliable communications? Are all of your friends on the same local area network? Do you really like shouting into the void? Well, Void Shouter is for you!

## How does it work?

It doesn’t, really. Maybe, kinda sorta.

## How is it supposed to work?

You type something, and it shouts a UDP multicast packet into the void.

## Who receives it?

Anybody who is connected to your local area network, i.e. your friends. Assuming you have any.

## But how do they receive it?

If they’re also running Void Shouter, they’ll receive it.

## But they’re not running Void Shouter!

Well, then I guess in that case they won’t receive it. But if they are, they will. Possibly, anyway.

## But I thought you said anybody connected to your local area network would receive it? You didn’t exactly qualify it.

Okay you got me, I wasn’t being specific, but at the same time I think you’re being too pedantic. If you want to shout about it, feel free to run:

```
% go build .
% ./voidshouter
<shout here>
```

## “Screenshot”

```text
% voidshouter 
2023/05/29 18:20:37 listening on [ff02::401d]:16413 
2023/05/29 18:20:42 fe80::9bf6:0e59:c196:e736 <alice> watson, come here, I want to see you
2023/05/29 18:20:44 <bob> who the heck is watson?
2023/05/29 18:20:47 <bob> alice?
2023/05/29 18:20:48 fe80::9bf6:0e59:c196:e736 <alice> urghghghkkkhkhkkh...
<bob> I appreciate the onomatopoeia, but█
```

## Support

If you need help, run `voidshouter`. If there is anybody else on your local area network, they might help out.

## Premium Support

Prefix each message with “I am very important, please pay full attention:”. That always works.

## Enterprise Support

Our people will talk to your people. You can speed up the process by telling our people where to get those giant novelty cheques printed.

Since sending messages over a network is a relatively novel concept, I fully expect this to be in high demand. Any moment now, the value will skyrocket, so the best thing you can do to prepare yourself for the imminent utopian future is signing a fixed term support contract (10 years at minimum).