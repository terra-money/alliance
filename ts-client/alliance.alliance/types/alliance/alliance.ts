/* eslint-disable */
import { Timestamp } from "../google/protobuf/timestamp";
import { RewardHistory } from "../alliance/params";
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "alliance.alliance";

/** key: denom value: AllianceAsset */
export interface AllianceAsset {
  /** Denom of the asset. It could either be a native token or an IBC token */
  denom: string;
  /**
   * The reward weight specifies the ratio of rewards that will be given to each alliance asset
   * It does not need to sum to 1. rate = weight / total_weight
   * Native asset is always assumed to have a weight of 1.
   */
  rewardWeight: string;
  /**
   * A positive take rate is used for liquid staking derivatives. It defines an annualized reward rate that
   * will be redirected to the distribution rewards pool
   */
  takeRate: string;
  totalTokens: string;
  totalValidatorShares: string;
  rewardStartTime: Date | undefined;
}

export interface AddAssetProposal {
  title: string;
  description: string;
  asset: AllianceAsset | undefined;
}

export interface RemoveAssetProposal {
  title: string;
  description: string;
  denom: string;
}

export interface UpdateAssetProposal {
  title: string;
  description: string;
  asset: AllianceAsset | undefined;
}

export interface QueuedRewardRateChange {
  denom: string;
  prevRewardRate: string;
}

export interface RewardRateChangeSnapshot {
  prevRewardWeight: string;
  rewardHistories: RewardHistory[];
}

const baseAllianceAsset: object = {
  denom: "",
  rewardWeight: "",
  takeRate: "",
  totalTokens: "",
  totalValidatorShares: "",
};

export const AllianceAsset = {
  encode(message: AllianceAsset, writer: Writer = Writer.create()): Writer {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }
    if (message.rewardWeight !== "") {
      writer.uint32(18).string(message.rewardWeight);
    }
    if (message.takeRate !== "") {
      writer.uint32(26).string(message.takeRate);
    }
    if (message.totalTokens !== "") {
      writer.uint32(34).string(message.totalTokens);
    }
    if (message.totalValidatorShares !== "") {
      writer.uint32(42).string(message.totalValidatorShares);
    }
    if (message.rewardStartTime !== undefined) {
      Timestamp.encode(
        toTimestamp(message.rewardStartTime),
        writer.uint32(50).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): AllianceAsset {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseAllianceAsset } as AllianceAsset;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.denom = reader.string();
          break;
        case 2:
          message.rewardWeight = reader.string();
          break;
        case 3:
          message.takeRate = reader.string();
          break;
        case 4:
          message.totalTokens = reader.string();
          break;
        case 5:
          message.totalValidatorShares = reader.string();
          break;
        case 6:
          message.rewardStartTime = fromTimestamp(
            Timestamp.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AllianceAsset {
    const message = { ...baseAllianceAsset } as AllianceAsset;
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = String(object.denom);
    } else {
      message.denom = "";
    }
    if (object.rewardWeight !== undefined && object.rewardWeight !== null) {
      message.rewardWeight = String(object.rewardWeight);
    } else {
      message.rewardWeight = "";
    }
    if (object.takeRate !== undefined && object.takeRate !== null) {
      message.takeRate = String(object.takeRate);
    } else {
      message.takeRate = "";
    }
    if (object.totalTokens !== undefined && object.totalTokens !== null) {
      message.totalTokens = String(object.totalTokens);
    } else {
      message.totalTokens = "";
    }
    if (
      object.totalValidatorShares !== undefined &&
      object.totalValidatorShares !== null
    ) {
      message.totalValidatorShares = String(object.totalValidatorShares);
    } else {
      message.totalValidatorShares = "";
    }
    if (
      object.rewardStartTime !== undefined &&
      object.rewardStartTime !== null
    ) {
      message.rewardStartTime = fromJsonTimestamp(object.rewardStartTime);
    } else {
      message.rewardStartTime = undefined;
    }
    return message;
  },

  toJSON(message: AllianceAsset): unknown {
    const obj: any = {};
    message.denom !== undefined && (obj.denom = message.denom);
    message.rewardWeight !== undefined &&
      (obj.rewardWeight = message.rewardWeight);
    message.takeRate !== undefined && (obj.takeRate = message.takeRate);
    message.totalTokens !== undefined &&
      (obj.totalTokens = message.totalTokens);
    message.totalValidatorShares !== undefined &&
      (obj.totalValidatorShares = message.totalValidatorShares);
    message.rewardStartTime !== undefined &&
      (obj.rewardStartTime =
        message.rewardStartTime !== undefined
          ? message.rewardStartTime.toISOString()
          : null);
    return obj;
  },

  fromPartial(object: DeepPartial<AllianceAsset>): AllianceAsset {
    const message = { ...baseAllianceAsset } as AllianceAsset;
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = object.denom;
    } else {
      message.denom = "";
    }
    if (object.rewardWeight !== undefined && object.rewardWeight !== null) {
      message.rewardWeight = object.rewardWeight;
    } else {
      message.rewardWeight = "";
    }
    if (object.takeRate !== undefined && object.takeRate !== null) {
      message.takeRate = object.takeRate;
    } else {
      message.takeRate = "";
    }
    if (object.totalTokens !== undefined && object.totalTokens !== null) {
      message.totalTokens = object.totalTokens;
    } else {
      message.totalTokens = "";
    }
    if (
      object.totalValidatorShares !== undefined &&
      object.totalValidatorShares !== null
    ) {
      message.totalValidatorShares = object.totalValidatorShares;
    } else {
      message.totalValidatorShares = "";
    }
    if (
      object.rewardStartTime !== undefined &&
      object.rewardStartTime !== null
    ) {
      message.rewardStartTime = object.rewardStartTime;
    } else {
      message.rewardStartTime = undefined;
    }
    return message;
  },
};

const baseAddAssetProposal: object = { title: "", description: "" };

export const AddAssetProposal = {
  encode(message: AddAssetProposal, writer: Writer = Writer.create()): Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.asset !== undefined) {
      AllianceAsset.encode(message.asset, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): AddAssetProposal {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseAddAssetProposal } as AddAssetProposal;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.asset = AllianceAsset.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddAssetProposal {
    const message = { ...baseAddAssetProposal } as AddAssetProposal;
    if (object.title !== undefined && object.title !== null) {
      message.title = String(object.title);
    } else {
      message.title = "";
    }
    if (object.description !== undefined && object.description !== null) {
      message.description = String(object.description);
    } else {
      message.description = "";
    }
    if (object.asset !== undefined && object.asset !== null) {
      message.asset = AllianceAsset.fromJSON(object.asset);
    } else {
      message.asset = undefined;
    }
    return message;
  },

  toJSON(message: AddAssetProposal): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined &&
      (obj.description = message.description);
    message.asset !== undefined &&
      (obj.asset = message.asset
        ? AllianceAsset.toJSON(message.asset)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<AddAssetProposal>): AddAssetProposal {
    const message = { ...baseAddAssetProposal } as AddAssetProposal;
    if (object.title !== undefined && object.title !== null) {
      message.title = object.title;
    } else {
      message.title = "";
    }
    if (object.description !== undefined && object.description !== null) {
      message.description = object.description;
    } else {
      message.description = "";
    }
    if (object.asset !== undefined && object.asset !== null) {
      message.asset = AllianceAsset.fromPartial(object.asset);
    } else {
      message.asset = undefined;
    }
    return message;
  },
};

const baseRemoveAssetProposal: object = {
  title: "",
  description: "",
  denom: "",
};

export const RemoveAssetProposal = {
  encode(
    message: RemoveAssetProposal,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.denom !== "") {
      writer.uint32(26).string(message.denom);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): RemoveAssetProposal {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseRemoveAssetProposal } as RemoveAssetProposal;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.denom = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveAssetProposal {
    const message = { ...baseRemoveAssetProposal } as RemoveAssetProposal;
    if (object.title !== undefined && object.title !== null) {
      message.title = String(object.title);
    } else {
      message.title = "";
    }
    if (object.description !== undefined && object.description !== null) {
      message.description = String(object.description);
    } else {
      message.description = "";
    }
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = String(object.denom);
    } else {
      message.denom = "";
    }
    return message;
  },

  toJSON(message: RemoveAssetProposal): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined &&
      (obj.description = message.description);
    message.denom !== undefined && (obj.denom = message.denom);
    return obj;
  },

  fromPartial(object: DeepPartial<RemoveAssetProposal>): RemoveAssetProposal {
    const message = { ...baseRemoveAssetProposal } as RemoveAssetProposal;
    if (object.title !== undefined && object.title !== null) {
      message.title = object.title;
    } else {
      message.title = "";
    }
    if (object.description !== undefined && object.description !== null) {
      message.description = object.description;
    } else {
      message.description = "";
    }
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = object.denom;
    } else {
      message.denom = "";
    }
    return message;
  },
};

const baseUpdateAssetProposal: object = { title: "", description: "" };

export const UpdateAssetProposal = {
  encode(
    message: UpdateAssetProposal,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.asset !== undefined) {
      AllianceAsset.encode(message.asset, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): UpdateAssetProposal {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseUpdateAssetProposal } as UpdateAssetProposal;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.asset = AllianceAsset.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateAssetProposal {
    const message = { ...baseUpdateAssetProposal } as UpdateAssetProposal;
    if (object.title !== undefined && object.title !== null) {
      message.title = String(object.title);
    } else {
      message.title = "";
    }
    if (object.description !== undefined && object.description !== null) {
      message.description = String(object.description);
    } else {
      message.description = "";
    }
    if (object.asset !== undefined && object.asset !== null) {
      message.asset = AllianceAsset.fromJSON(object.asset);
    } else {
      message.asset = undefined;
    }
    return message;
  },

  toJSON(message: UpdateAssetProposal): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined &&
      (obj.description = message.description);
    message.asset !== undefined &&
      (obj.asset = message.asset
        ? AllianceAsset.toJSON(message.asset)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<UpdateAssetProposal>): UpdateAssetProposal {
    const message = { ...baseUpdateAssetProposal } as UpdateAssetProposal;
    if (object.title !== undefined && object.title !== null) {
      message.title = object.title;
    } else {
      message.title = "";
    }
    if (object.description !== undefined && object.description !== null) {
      message.description = object.description;
    } else {
      message.description = "";
    }
    if (object.asset !== undefined && object.asset !== null) {
      message.asset = AllianceAsset.fromPartial(object.asset);
    } else {
      message.asset = undefined;
    }
    return message;
  },
};

const baseQueuedRewardRateChange: object = { denom: "", prevRewardRate: "" };

export const QueuedRewardRateChange = {
  encode(
    message: QueuedRewardRateChange,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }
    if (message.prevRewardRate !== "") {
      writer.uint32(18).string(message.prevRewardRate);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueuedRewardRateChange {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueuedRewardRateChange } as QueuedRewardRateChange;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.denom = reader.string();
          break;
        case 2:
          message.prevRewardRate = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueuedRewardRateChange {
    const message = { ...baseQueuedRewardRateChange } as QueuedRewardRateChange;
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = String(object.denom);
    } else {
      message.denom = "";
    }
    if (object.prevRewardRate !== undefined && object.prevRewardRate !== null) {
      message.prevRewardRate = String(object.prevRewardRate);
    } else {
      message.prevRewardRate = "";
    }
    return message;
  },

  toJSON(message: QueuedRewardRateChange): unknown {
    const obj: any = {};
    message.denom !== undefined && (obj.denom = message.denom);
    message.prevRewardRate !== undefined &&
      (obj.prevRewardRate = message.prevRewardRate);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueuedRewardRateChange>
  ): QueuedRewardRateChange {
    const message = { ...baseQueuedRewardRateChange } as QueuedRewardRateChange;
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = object.denom;
    } else {
      message.denom = "";
    }
    if (object.prevRewardRate !== undefined && object.prevRewardRate !== null) {
      message.prevRewardRate = object.prevRewardRate;
    } else {
      message.prevRewardRate = "";
    }
    return message;
  },
};

const baseRewardRateChangeSnapshot: object = { prevRewardWeight: "" };

export const RewardRateChangeSnapshot = {
  encode(
    message: RewardRateChangeSnapshot,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.prevRewardWeight !== "") {
      writer.uint32(10).string(message.prevRewardWeight);
    }
    for (const v of message.rewardHistories) {
      RewardHistory.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): RewardRateChangeSnapshot {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseRewardRateChangeSnapshot,
    } as RewardRateChangeSnapshot;
    message.rewardHistories = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.prevRewardWeight = reader.string();
          break;
        case 2:
          message.rewardHistories.push(
            RewardHistory.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RewardRateChangeSnapshot {
    const message = {
      ...baseRewardRateChangeSnapshot,
    } as RewardRateChangeSnapshot;
    message.rewardHistories = [];
    if (
      object.prevRewardWeight !== undefined &&
      object.prevRewardWeight !== null
    ) {
      message.prevRewardWeight = String(object.prevRewardWeight);
    } else {
      message.prevRewardWeight = "";
    }
    if (
      object.rewardHistories !== undefined &&
      object.rewardHistories !== null
    ) {
      for (const e of object.rewardHistories) {
        message.rewardHistories.push(RewardHistory.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: RewardRateChangeSnapshot): unknown {
    const obj: any = {};
    message.prevRewardWeight !== undefined &&
      (obj.prevRewardWeight = message.prevRewardWeight);
    if (message.rewardHistories) {
      obj.rewardHistories = message.rewardHistories.map((e) =>
        e ? RewardHistory.toJSON(e) : undefined
      );
    } else {
      obj.rewardHistories = [];
    }
    return obj;
  },

  fromPartial(
    object: DeepPartial<RewardRateChangeSnapshot>
  ): RewardRateChangeSnapshot {
    const message = {
      ...baseRewardRateChangeSnapshot,
    } as RewardRateChangeSnapshot;
    message.rewardHistories = [];
    if (
      object.prevRewardWeight !== undefined &&
      object.prevRewardWeight !== null
    ) {
      message.prevRewardWeight = object.prevRewardWeight;
    } else {
      message.prevRewardWeight = "";
    }
    if (
      object.rewardHistories !== undefined &&
      object.rewardHistories !== null
    ) {
      for (const e of object.rewardHistories) {
        message.rewardHistories.push(RewardHistory.fromPartial(e));
      }
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
