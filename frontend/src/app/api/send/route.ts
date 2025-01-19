import { NextRequest, NextResponse } from "next/server";
import z from "zod";

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

    const res: string = name + " " + key;

    return NextResponse.json(res, {
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
