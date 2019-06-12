# Presentation

* [Microservice Architectures: What They Are and Why You Should Use Them](https://blog.newrelic.com/technology/microservices-what-they-are-why-to-use-them/)

## Introduction & Background

Good morning, class! Our group consists of Jams Cao, Lin Xu Yalun and I. Today our topic is Microservices. 

Microservices is a hot topic in software development circles these days. And for some very good reasons. Now let’s tell you What They Are and Why You Should Use Them.

### Monoliths: the way we were

First let's look at the traditional approach—Monolithic approach. Here is a diagram of how many monolithic applications start out.: *(a graph in slides)* 

Traditionally, software developers created large, monolithic applications. A single monolith would contain all the code for all the business activities an application performed. And what happens if your application turns out to be successful? Users like it and begin to depend on it. Traffic increases dramatically. And almost inevitably, users request improvements and additional features, so more developers are roped in to work on the growing application. Before too long, what will your application look like?

Just like this: *(a graph in slides)* Your once-simple application has become large and complex. Multiple independent development teams are working on it at the same time. 

Fortunately, this scenario is not inevitable. You can rebuild and re-architect your applications to scale with your company’s needs, not against them. Using a modern, microservices-based application architecture is an increasingly popular technique for building applications that can scale without miring your organization in monolithic muck.

## Body

### Overview

So what exactly is microservice architecture? The main idea behind microservice architecture is that applications are simpler to build and maintain when broken down into smaller pieces that work seamlessly together.

Your application looks more like this: *a graph in slides*

This diagram shows an application constructed as a series of microservices. As you can see, microservice architectures let you split applications into distinct independent services, and each managed by individual groups in paralle. 

By delegating the responsibility to independent groups, a difficult and complicated application become simple and manageable. One group can work independently without impacting the work of other developers in other groups working on the same overall application.

In short, your application can grow as your company and its requirements grow easily. 

### Microservices and containers

It’s difficult to talk about microservices without also talking about containers. Within a container, an application is deployed. A container offers a siginificant improvement over the VM, which must have an OS and could only manage one application. 

In the past, people usually deploy their applications in VMs, like the graph on the left. As you can see, one VM has one OS and one application. The OS usually costs many resources, which is, however, a heavy burden to an application. However, developers had no other choices in the past. And here comes the container, invented in 2013. Like the graph on the right, containers don't have an OS,  conainters are above OS. So the amazing thing is that you can set up hundreds of containers on an OS, which means you can launch hundres of applications in a VM rather than just one application before. Remember the architecture of microservice? We have a lot of small applications working together. And that's why many developers feel that containers are a natural fit for microservices.

### A Survey of Microservices(leanIX 2017)

To present microservices applied in the real world, here we show a survey that was carried out by leanIX in 2017. (动画) This pie chart shows the frequency of application releases of companies that do not adopt microservices. As we can see, most companies only release multiple times a year, and very few companies release multiple times per week. However, things become different with microservices applying. (动画) From this new chart, we could find that most companies have multiple times of release per week. Actually, the more frequently the company release applications, the more rapidly it could react to changing requirements of customers. This could decide on the market share of applications significantly, especially of applications in developing markets. (动画) For example, Wechat released several times per day in its early time. Thus, if you want to build a super application like Wechat, you’d better adopt microservices to release more frequently.  (切换ppt) 

The bar chart in the left shows the difficulties companies meet when adopting microservices. Missing knowledge and people, legacy processes and higher complexity are all common hurdles. (动画)
 The complex microservices architecture, as the graph displayed in the right, could be blamed for these hurdles. There are often more components involved in microservices architecture, and those components have more interconnections. 

## Conclusion
In conclusion, developing applications in microservices would improve scalability, maintainability, and agility. Although there are hurdles due to the complexity of the architecture, the benefits worth it. (动画)
 As the graph shows, up to 71% of companies would intensify the usage of microservices. Therefore,microservices is strongly suggested to be adopted in your applications, especially for newly set up ones. That’s it for our talk today and thank you for your attention. Any questions, please?

