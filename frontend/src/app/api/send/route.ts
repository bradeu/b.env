import { NextRequest, NextResponse } from "next/server";
import z from "zod";
import amqp from "amqplib";

const schema = z.object({
  name: z.string().min(1),
  key: z.string().min(1),
});

export async function POST(req: NextRequest) {
  try {
    const body = await req.json();

    // just some string validation with zod
    const { name, key } = schema.parse(body);

    // encrypt data here

    // api call here
    const exchange = "intercept_exchange";
    const queue = "intercept";
    const route = "mail_route";

    const uri = "localhost";
    const user = "segal";
    const pass = "xxxxxxxx";
    const url = `amqp://${user}:${pass}@${uri}`;

    const connection = await amqp.connect(url);
    const channel = await connection.createChannel();

    await channel.assertExchange(exchange, "direct");
    await channel.assertQueue("mail");
    await channel.bindQueue(queue, exchange, route);

    const sent = channel.publish(
      exchange,
      route,
      Buffer.from(JSON.stringify({ name, key }))
    );

    if (sent)
      console.info(
        `${name} - Sent message to ${exchange} -> ${route} ${JSON.stringify({
          name,
          key,
        })}`
      );

    return NextResponse.json("", {
      status: 302,
      headers: {
        Location: "/",
      },
    });
  } catch {
    return NextResponse.json(
      { error: "Something Failed" },
      {
        status: 500,
      }
    );
  }
}
