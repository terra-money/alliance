/* eslint-disable */
import { Timestamp } from "../google/protobuf/timestamp";
import { Duration } from "../google/protobuf/duration";
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "alliance.alliance";

export interface Params {
  rewardDelayTime: Duration | undefined;
  /** Time interval between consecutive applications of `take_rate` */
  rewardClaimInterval: Duration | undefined;
  /** Last application of `take_rate` on assets */
  lastRewardClaimTime: Date | undefined;
  globalRewardIndices: RewardIndex[];
}

export interface RewardIndex {
  denom: string;
  index: string;
}

const baseParams: object = {};

export const Params = {
  encode(message: Params, writer: Writer = Writer.create()): Writer {
    if (message.rewardDelayTime !== undefined) {
      Duration.encode(
        message.rewardDelayTime,
        writer.uint32(10).fork()
      ).ldelim();
    }
    if (message.rewardClaimInterval !== undefined) {
      Duration.encode(
        message.rewardClaimInterval,
        writer.uint32(18).fork()
      ).ldelim();
    }
    if (message.lastRewardClaimTime !== undefined) {
      Timestamp.encode(
        toTimestamp(message.lastRewardClaimTime),
        writer.uint32(26).fork()
      ).ldelim();
    }
    for (const v of message.globalRewardIndices) {
      RewardIndex.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): Params {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseParams } as Params;
    message.globalRewardIndices = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rewardDelayTime = Duration.decode(reader, reader.uint32());
          break;
        case 2:
          message.rewardClaimInterval = Duration.decode(
            reader,
            reader.uint32()
          );
          break;
        case 3:
          message.lastRewardClaimTime = fromTimestamp(
            Timestamp.decode(reader, reader.uint32())
          );
          break;
        case 4:
          message.globalRewardIndices.push(
            RewardIndex.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Params {
    const message = { ...baseParams } as Params;
    message.globalRewardIndices = [];
    if (
      object.rewardDelayTime !== undefined &&
      object.rewardDelayTime !== null
    ) {
      message.rewardDelayTime = Duration.fromJSON(object.rewardDelayTime);
    } else {
      message.rewardDelayTime = undefined;
    }
    if (
      object.rewardClaimInterval !== undefined &&
      object.rewardClaimInterval !== null
    ) {
      message.rewardClaimInterval = Duration.fromJSON(
        object.rewardClaimInterval
      );
    } else {
      message.rewardClaimInterval = undefined;
    }
    if (
      object.lastRewardClaimTime !== undefined &&
      object.lastRewardClaimTime !== null
    ) {
      message.lastRewardClaimTime = fromJsonTimestamp(
        object.lastRewardClaimTime
      );
    } else {
      message.lastRewardClaimTime = undefined;
    }
    if (
      object.globalRewardIndices !== undefined &&
      object.globalRewardIndices !== null
    ) {
      for (const e of object.globalRewardIndices) {
        message.globalRewardIndices.push(RewardIndex.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: Params): unknown {
    const obj: any = {};
    message.rewardDelayTime !== undefined &&
      (obj.rewardDelayTime = message.rewardDelayTime
        ? Duration.toJSON(message.rewardDelayTime)
        : undefined);
    message.rewardClaimInterval !== undefined &&
      (obj.rewardClaimInterval = message.rewardClaimInterval
        ? Duration.toJSON(message.rewardClaimInterval)
        : undefined);
    message.lastRewardClaimTime !== undefined &&
      (obj.lastRewardClaimTime =
        message.lastRewardClaimTime !== undefined
          ? message.lastRewardClaimTime.toISOString()
          : null);
    if (message.globalRewardIndices) {
      obj.globalRewardIndices = message.globalRewardIndices.map((e) =>
        e ? RewardIndex.toJSON(e) : undefined
      );
    } else {
      obj.globalRewardIndices = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<Params>): Params {
    const message = { ...baseParams } as Params;
    message.globalRewardIndices = [];
    if (
      object.rewardDelayTime !== undefined &&
      object.rewardDelayTime !== null
    ) {
      message.rewardDelayTime = Duration.fromPartial(object.rewardDelayTime);
    } else {
      message.rewardDelayTime = undefined;
    }
    if (
      object.rewardClaimInterval !== undefined &&
      object.rewardClaimInterval !== null
    ) {
      message.rewardClaimInterval = Duration.fromPartial(
        object.rewardClaimInterval
      );
    } else {
      message.rewardClaimInterval = undefined;
    }
    if (
      object.lastRewardClaimTime !== undefined &&
      object.lastRewardClaimTime !== null
    ) {
      message.lastRewardClaimTime = object.lastRewardClaimTime;
    } else {
      message.lastRewardClaimTime = undefined;
    }
    if (
      object.globalRewardIndices !== undefined &&
      object.globalRewardIndices !== null
    ) {
      for (const e of object.globalRewardIndices) {
        message.globalRewardIndices.push(RewardIndex.fromPartial(e));
      }
    }
    return message;
  },
};

const baseRewardIndex: object = { denom: "", index: "" };

export const RewardIndex = {
  encode(message: RewardIndex, writer: Writer = Writer.create()): Writer {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }
    if (message.index !== "") {
      writer.uint32(18).string(message.index);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): RewardIndex {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseRewardIndex } as RewardIndex;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.denom = reader.string();
          break;
        case 2:
          message.index = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RewardIndex {
    const message = { ...baseRewardIndex } as RewardIndex;
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = String(object.denom);
    } else {
      message.denom = "";
    }
    if (object.index !== undefined && object.index !== null) {
      message.index = String(object.index);
    } else {
      message.index = "";
    }
    return message;
  },

  toJSON(message: RewardIndex): unknown {
    const obj: any = {};
    message.denom !== undefined && (obj.denom = message.denom);
    message.index !== undefined && (obj.index = message.index);
    return obj;
  },

  fromPartial(object: DeepPartial<RewardIndex>): RewardIndex {
    const message = { ...baseRewardIndex } as RewardIndex;
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = object.denom;
    } else {
      message.denom = "";
    }
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    } else {
      message.index = "";
    }
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}
