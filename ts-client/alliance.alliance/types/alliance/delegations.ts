/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { RewardHistory } from "../alliance/params";
import { Coin, DecCoin } from "../cosmos/base/v1beta1/coin";

export const protobufPackage = "alliance.alliance";

export interface Delegation {
  /** delegator_address is the bech32-encoded address of the delegator. */
  delegatorAddress: string;
  /** validator_address is the bech32-encoded address of the validator. */
  validatorAddress: string;
  /** denom of token staked */
  denom: string;
  /** shares define the delegation shares received. */
  shares: string;
  rewardHistory: RewardHistory[];
  lastRewardClaimHeight: number;
}

/**
 * DelegationResponse is equivalent to Delegation except that it contains a
 * balance in addition to shares which is more suitable for client responses.
 */
export interface DelegationResponse {
  delegation: Delegation | undefined;
  balance: Coin | undefined;
}

export interface Redelegation {
  delegatorAddress: string;
  srcValidatorAddress: string;
  dstValidatorAddress: string;
  balance: Coin | undefined;
}

export interface QueuedRedelegation {
  entries: Redelegation[];
}

export interface Undelegation {
  delegatorAddress: string;
  validatorAddress: string;
  balance: Coin | undefined;
}

export interface QueuedUndelegation {
  entries: Undelegation[];
}

export interface AllianceValidatorInfo {
  globalRewardHistory: RewardHistory[];
  totalDelegatorShares: DecCoin[];
  validatorShares: DecCoin[];
}

const baseDelegation: object = {
  delegatorAddress: "",
  validatorAddress: "",
  denom: "",
  shares: "",
  lastRewardClaimHeight: 0,
};

export const Delegation = {
  encode(message: Delegation, writer: Writer = Writer.create()): Writer {
    if (message.delegatorAddress !== "") {
      writer.uint32(10).string(message.delegatorAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    if (message.denom !== "") {
      writer.uint32(26).string(message.denom);
    }
    if (message.shares !== "") {
      writer.uint32(34).string(message.shares);
    }
    for (const v of message.rewardHistory) {
      RewardHistory.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.lastRewardClaimHeight !== 0) {
      writer.uint32(48).uint64(message.lastRewardClaimHeight);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): Delegation {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseDelegation } as Delegation;
    message.rewardHistory = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.delegatorAddress = reader.string();
          break;
        case 2:
          message.validatorAddress = reader.string();
          break;
        case 3:
          message.denom = reader.string();
          break;
        case 4:
          message.shares = reader.string();
          break;
        case 5:
          message.rewardHistory.push(
            RewardHistory.decode(reader, reader.uint32())
          );
          break;
        case 6:
          message.lastRewardClaimHeight = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Delegation {
    const message = { ...baseDelegation } as Delegation;
    message.rewardHistory = [];
    if (
      object.delegatorAddress !== undefined &&
      object.delegatorAddress !== null
    ) {
      message.delegatorAddress = String(object.delegatorAddress);
    } else {
      message.delegatorAddress = "";
    }
    if (
      object.validatorAddress !== undefined &&
      object.validatorAddress !== null
    ) {
      message.validatorAddress = String(object.validatorAddress);
    } else {
      message.validatorAddress = "";
    }
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = String(object.denom);
    } else {
      message.denom = "";
    }
    if (object.shares !== undefined && object.shares !== null) {
      message.shares = String(object.shares);
    } else {
      message.shares = "";
    }
    if (object.rewardHistory !== undefined && object.rewardHistory !== null) {
      for (const e of object.rewardHistory) {
        message.rewardHistory.push(RewardHistory.fromJSON(e));
      }
    }
    if (
      object.lastRewardClaimHeight !== undefined &&
      object.lastRewardClaimHeight !== null
    ) {
      message.lastRewardClaimHeight = Number(object.lastRewardClaimHeight);
    } else {
      message.lastRewardClaimHeight = 0;
    }
    return message;
  },

  toJSON(message: Delegation): unknown {
    const obj: any = {};
    message.delegatorAddress !== undefined &&
      (obj.delegatorAddress = message.delegatorAddress);
    message.validatorAddress !== undefined &&
      (obj.validatorAddress = message.validatorAddress);
    message.denom !== undefined && (obj.denom = message.denom);
    message.shares !== undefined && (obj.shares = message.shares);
    if (message.rewardHistory) {
      obj.rewardHistory = message.rewardHistory.map((e) =>
        e ? RewardHistory.toJSON(e) : undefined
      );
    } else {
      obj.rewardHistory = [];
    }
    message.lastRewardClaimHeight !== undefined &&
      (obj.lastRewardClaimHeight = message.lastRewardClaimHeight);
    return obj;
  },

  fromPartial(object: DeepPartial<Delegation>): Delegation {
    const message = { ...baseDelegation } as Delegation;
    message.rewardHistory = [];
    if (
      object.delegatorAddress !== undefined &&
      object.delegatorAddress !== null
    ) {
      message.delegatorAddress = object.delegatorAddress;
    } else {
      message.delegatorAddress = "";
    }
    if (
      object.validatorAddress !== undefined &&
      object.validatorAddress !== null
    ) {
      message.validatorAddress = object.validatorAddress;
    } else {
      message.validatorAddress = "";
    }
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = object.denom;
    } else {
      message.denom = "";
    }
    if (object.shares !== undefined && object.shares !== null) {
      message.shares = object.shares;
    } else {
      message.shares = "";
    }
    if (object.rewardHistory !== undefined && object.rewardHistory !== null) {
      for (const e of object.rewardHistory) {
        message.rewardHistory.push(RewardHistory.fromPartial(e));
      }
    }
    if (
      object.lastRewardClaimHeight !== undefined &&
      object.lastRewardClaimHeight !== null
    ) {
      message.lastRewardClaimHeight = object.lastRewardClaimHeight;
    } else {
      message.lastRewardClaimHeight = 0;
    }
    return message;
  },
};

const baseDelegationResponse: object = {};

export const DelegationResponse = {
  encode(
    message: DelegationResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.delegation !== undefined) {
      Delegation.encode(message.delegation, writer.uint32(10).fork()).ldelim();
    }
    if (message.balance !== undefined) {
      Coin.encode(message.balance, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): DelegationResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseDelegationResponse } as DelegationResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.delegation = Delegation.decode(reader, reader.uint32());
          break;
        case 2:
          message.balance = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DelegationResponse {
    const message = { ...baseDelegationResponse } as DelegationResponse;
    if (object.delegation !== undefined && object.delegation !== null) {
      message.delegation = Delegation.fromJSON(object.delegation);
    } else {
      message.delegation = undefined;
    }
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = Coin.fromJSON(object.balance);
    } else {
      message.balance = undefined;
    }
    return message;
  },

  toJSON(message: DelegationResponse): unknown {
    const obj: any = {};
    message.delegation !== undefined &&
      (obj.delegation = message.delegation
        ? Delegation.toJSON(message.delegation)
        : undefined);
    message.balance !== undefined &&
      (obj.balance = message.balance
        ? Coin.toJSON(message.balance)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<DelegationResponse>): DelegationResponse {
    const message = { ...baseDelegationResponse } as DelegationResponse;
    if (object.delegation !== undefined && object.delegation !== null) {
      message.delegation = Delegation.fromPartial(object.delegation);
    } else {
      message.delegation = undefined;
    }
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = Coin.fromPartial(object.balance);
    } else {
      message.balance = undefined;
    }
    return message;
  },
};

const baseRedelegation: object = {
  delegatorAddress: "",
  srcValidatorAddress: "",
  dstValidatorAddress: "",
};

export const Redelegation = {
  encode(message: Redelegation, writer: Writer = Writer.create()): Writer {
    if (message.delegatorAddress !== "") {
      writer.uint32(10).string(message.delegatorAddress);
    }
    if (message.srcValidatorAddress !== "") {
      writer.uint32(18).string(message.srcValidatorAddress);
    }
    if (message.dstValidatorAddress !== "") {
      writer.uint32(26).string(message.dstValidatorAddress);
    }
    if (message.balance !== undefined) {
      Coin.encode(message.balance, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): Redelegation {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseRedelegation } as Redelegation;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.delegatorAddress = reader.string();
          break;
        case 2:
          message.srcValidatorAddress = reader.string();
          break;
        case 3:
          message.dstValidatorAddress = reader.string();
          break;
        case 4:
          message.balance = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Redelegation {
    const message = { ...baseRedelegation } as Redelegation;
    if (
      object.delegatorAddress !== undefined &&
      object.delegatorAddress !== null
    ) {
      message.delegatorAddress = String(object.delegatorAddress);
    } else {
      message.delegatorAddress = "";
    }
    if (
      object.srcValidatorAddress !== undefined &&
      object.srcValidatorAddress !== null
    ) {
      message.srcValidatorAddress = String(object.srcValidatorAddress);
    } else {
      message.srcValidatorAddress = "";
    }
    if (
      object.dstValidatorAddress !== undefined &&
      object.dstValidatorAddress !== null
    ) {
      message.dstValidatorAddress = String(object.dstValidatorAddress);
    } else {
      message.dstValidatorAddress = "";
    }
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = Coin.fromJSON(object.balance);
    } else {
      message.balance = undefined;
    }
    return message;
  },

  toJSON(message: Redelegation): unknown {
    const obj: any = {};
    message.delegatorAddress !== undefined &&
      (obj.delegatorAddress = message.delegatorAddress);
    message.srcValidatorAddress !== undefined &&
      (obj.srcValidatorAddress = message.srcValidatorAddress);
    message.dstValidatorAddress !== undefined &&
      (obj.dstValidatorAddress = message.dstValidatorAddress);
    message.balance !== undefined &&
      (obj.balance = message.balance
        ? Coin.toJSON(message.balance)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<Redelegation>): Redelegation {
    const message = { ...baseRedelegation } as Redelegation;
    if (
      object.delegatorAddress !== undefined &&
      object.delegatorAddress !== null
    ) {
      message.delegatorAddress = object.delegatorAddress;
    } else {
      message.delegatorAddress = "";
    }
    if (
      object.srcValidatorAddress !== undefined &&
      object.srcValidatorAddress !== null
    ) {
      message.srcValidatorAddress = object.srcValidatorAddress;
    } else {
      message.srcValidatorAddress = "";
    }
    if (
      object.dstValidatorAddress !== undefined &&
      object.dstValidatorAddress !== null
    ) {
      message.dstValidatorAddress = object.dstValidatorAddress;
    } else {
      message.dstValidatorAddress = "";
    }
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = Coin.fromPartial(object.balance);
    } else {
      message.balance = undefined;
    }
    return message;
  },
};

const baseQueuedRedelegation: object = {};

export const QueuedRedelegation = {
  encode(
    message: QueuedRedelegation,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.entries) {
      Redelegation.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueuedRedelegation {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueuedRedelegation } as QueuedRedelegation;
    message.entries = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.entries.push(Redelegation.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueuedRedelegation {
    const message = { ...baseQueuedRedelegation } as QueuedRedelegation;
    message.entries = [];
    if (object.entries !== undefined && object.entries !== null) {
      for (const e of object.entries) {
        message.entries.push(Redelegation.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: QueuedRedelegation): unknown {
    const obj: any = {};
    if (message.entries) {
      obj.entries = message.entries.map((e) =>
        e ? Redelegation.toJSON(e) : undefined
      );
    } else {
      obj.entries = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<QueuedRedelegation>): QueuedRedelegation {
    const message = { ...baseQueuedRedelegation } as QueuedRedelegation;
    message.entries = [];
    if (object.entries !== undefined && object.entries !== null) {
      for (const e of object.entries) {
        message.entries.push(Redelegation.fromPartial(e));
      }
    }
    return message;
  },
};

const baseUndelegation: object = { delegatorAddress: "", validatorAddress: "" };

export const Undelegation = {
  encode(message: Undelegation, writer: Writer = Writer.create()): Writer {
    if (message.delegatorAddress !== "") {
      writer.uint32(10).string(message.delegatorAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    if (message.balance !== undefined) {
      Coin.encode(message.balance, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): Undelegation {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseUndelegation } as Undelegation;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.delegatorAddress = reader.string();
          break;
        case 2:
          message.validatorAddress = reader.string();
          break;
        case 3:
          message.balance = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Undelegation {
    const message = { ...baseUndelegation } as Undelegation;
    if (
      object.delegatorAddress !== undefined &&
      object.delegatorAddress !== null
    ) {
      message.delegatorAddress = String(object.delegatorAddress);
    } else {
      message.delegatorAddress = "";
    }
    if (
      object.validatorAddress !== undefined &&
      object.validatorAddress !== null
    ) {
      message.validatorAddress = String(object.validatorAddress);
    } else {
      message.validatorAddress = "";
    }
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = Coin.fromJSON(object.balance);
    } else {
      message.balance = undefined;
    }
    return message;
  },

  toJSON(message: Undelegation): unknown {
    const obj: any = {};
    message.delegatorAddress !== undefined &&
      (obj.delegatorAddress = message.delegatorAddress);
    message.validatorAddress !== undefined &&
      (obj.validatorAddress = message.validatorAddress);
    message.balance !== undefined &&
      (obj.balance = message.balance
        ? Coin.toJSON(message.balance)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<Undelegation>): Undelegation {
    const message = { ...baseUndelegation } as Undelegation;
    if (
      object.delegatorAddress !== undefined &&
      object.delegatorAddress !== null
    ) {
      message.delegatorAddress = object.delegatorAddress;
    } else {
      message.delegatorAddress = "";
    }
    if (
      object.validatorAddress !== undefined &&
      object.validatorAddress !== null
    ) {
      message.validatorAddress = object.validatorAddress;
    } else {
      message.validatorAddress = "";
    }
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = Coin.fromPartial(object.balance);
    } else {
      message.balance = undefined;
    }
    return message;
  },
};

const baseQueuedUndelegation: object = {};

export const QueuedUndelegation = {
  encode(
    message: QueuedUndelegation,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.entries) {
      Undelegation.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueuedUndelegation {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueuedUndelegation } as QueuedUndelegation;
    message.entries = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.entries.push(Undelegation.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueuedUndelegation {
    const message = { ...baseQueuedUndelegation } as QueuedUndelegation;
    message.entries = [];
    if (object.entries !== undefined && object.entries !== null) {
      for (const e of object.entries) {
        message.entries.push(Undelegation.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: QueuedUndelegation): unknown {
    const obj: any = {};
    if (message.entries) {
      obj.entries = message.entries.map((e) =>
        e ? Undelegation.toJSON(e) : undefined
      );
    } else {
      obj.entries = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<QueuedUndelegation>): QueuedUndelegation {
    const message = { ...baseQueuedUndelegation } as QueuedUndelegation;
    message.entries = [];
    if (object.entries !== undefined && object.entries !== null) {
      for (const e of object.entries) {
        message.entries.push(Undelegation.fromPartial(e));
      }
    }
    return message;
  },
};

const baseAllianceValidatorInfo: object = {};

export const AllianceValidatorInfo = {
  encode(
    message: AllianceValidatorInfo,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.globalRewardHistory) {
      RewardHistory.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.totalDelegatorShares) {
      DecCoin.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.validatorShares) {
      DecCoin.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): AllianceValidatorInfo {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseAllianceValidatorInfo } as AllianceValidatorInfo;
    message.globalRewardHistory = [];
    message.totalDelegatorShares = [];
    message.validatorShares = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.globalRewardHistory.push(
            RewardHistory.decode(reader, reader.uint32())
          );
          break;
        case 2:
          message.totalDelegatorShares.push(
            DecCoin.decode(reader, reader.uint32())
          );
          break;
        case 3:
          message.validatorShares.push(DecCoin.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AllianceValidatorInfo {
    const message = { ...baseAllianceValidatorInfo } as AllianceValidatorInfo;
    message.globalRewardHistory = [];
    message.totalDelegatorShares = [];
    message.validatorShares = [];
    if (
      object.globalRewardHistory !== undefined &&
      object.globalRewardHistory !== null
    ) {
      for (const e of object.globalRewardHistory) {
        message.globalRewardHistory.push(RewardHistory.fromJSON(e));
      }
    }
    if (
      object.totalDelegatorShares !== undefined &&
      object.totalDelegatorShares !== null
    ) {
      for (const e of object.totalDelegatorShares) {
        message.totalDelegatorShares.push(DecCoin.fromJSON(e));
      }
    }
    if (
      object.validatorShares !== undefined &&
      object.validatorShares !== null
    ) {
      for (const e of object.validatorShares) {
        message.validatorShares.push(DecCoin.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: AllianceValidatorInfo): unknown {
    const obj: any = {};
    if (message.globalRewardHistory) {
      obj.globalRewardHistory = message.globalRewardHistory.map((e) =>
        e ? RewardHistory.toJSON(e) : undefined
      );
    } else {
      obj.globalRewardHistory = [];
    }
    if (message.totalDelegatorShares) {
      obj.totalDelegatorShares = message.totalDelegatorShares.map((e) =>
        e ? DecCoin.toJSON(e) : undefined
      );
    } else {
      obj.totalDelegatorShares = [];
    }
    if (message.validatorShares) {
      obj.validatorShares = message.validatorShares.map((e) =>
        e ? DecCoin.toJSON(e) : undefined
      );
    } else {
      obj.validatorShares = [];
    }
    return obj;
  },

  fromPartial(
    object: DeepPartial<AllianceValidatorInfo>
  ): AllianceValidatorInfo {
    const message = { ...baseAllianceValidatorInfo } as AllianceValidatorInfo;
    message.globalRewardHistory = [];
    message.totalDelegatorShares = [];
    message.validatorShares = [];
    if (
      object.globalRewardHistory !== undefined &&
      object.globalRewardHistory !== null
    ) {
      for (const e of object.globalRewardHistory) {
        message.globalRewardHistory.push(RewardHistory.fromPartial(e));
      }
    }
    if (
      object.totalDelegatorShares !== undefined &&
      object.totalDelegatorShares !== null
    ) {
      for (const e of object.totalDelegatorShares) {
        message.totalDelegatorShares.push(DecCoin.fromPartial(e));
      }
    }
    if (
      object.validatorShares !== undefined &&
      object.validatorShares !== null
    ) {
      for (const e of object.validatorShares) {
        message.validatorShares.push(DecCoin.fromPartial(e));
      }
    }
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

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

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (util.Long !== Long) {
  util.Long = Long as any;
  configure();
}
