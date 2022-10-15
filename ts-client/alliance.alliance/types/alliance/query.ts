/* eslint-disable */
import { Reader, Writer } from "protobufjs/minimal";
import { Params } from "../alliance/params";
import {
  PageRequest,
  PageResponse,
} from "../cosmos/base/query/v1beta1/pagination";
import { AllianceAsset } from "../alliance/alliance";
import { Delegation } from "../alliance/delegations";

export const protobufPackage = "alliance.alliance";

export interface QueryParamsRequest {}

export interface QueryParamsResponse {
  params: Params | undefined;
}

export interface QueryAlliancesRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAlliancesResponse {
  assets: AllianceAsset[];
  pagination: PageResponse | undefined;
}

export interface QueryAllianceRequest {
  denom: string;
}

export interface QueryAllianceResponse {
  alliance: AllianceAsset | undefined;
}

export interface QueryAllianceDelegationsRequest {
  denom: string;
  pagination: PageResponse | undefined;
}

export interface QueryAllianceDelegationsResponse {
  delegations: Delegation[];
}

export interface QueryAllianceDelegationRequest {
  denom: string;
  delegatorAddress: string;
  pagination: PageResponse | undefined;
}

export interface QueryAllianceDelegationResponse {
  delegations: Delegation | undefined;
}

const baseQueryParamsRequest: object = {};

export const QueryParamsRequest = {
  encode(_: QueryParamsRequest, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryParamsRequest } as QueryParamsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): QueryParamsRequest {
    const message = { ...baseQueryParamsRequest } as QueryParamsRequest;
    return message;
  },

  toJSON(_: QueryParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<QueryParamsRequest>): QueryParamsRequest {
    const message = { ...baseQueryParamsRequest } as QueryParamsRequest;
    return message;
  },
};

const baseQueryParamsResponse: object = {};

export const QueryParamsResponse = {
  encode(
    message: QueryParamsResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryParamsResponse } as QueryParamsResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryParamsResponse {
    const message = { ...baseQueryParamsResponse } as QueryParamsResponse;
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromJSON(object.params);
    } else {
      message.params = undefined;
    }
    return message;
  },

  toJSON(message: QueryParamsResponse): unknown {
    const obj: any = {};
    message.params !== undefined &&
      (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<QueryParamsResponse>): QueryParamsResponse {
    const message = { ...baseQueryParamsResponse } as QueryParamsResponse;
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromPartial(object.params);
    } else {
      message.params = undefined;
    }
    return message;
  },
};

const baseQueryAlliancesRequest: object = {};

export const QueryAlliancesRequest = {
  encode(
    message: QueryAlliancesRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryAlliancesRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryAlliancesRequest } as QueryAlliancesRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAlliancesRequest {
    const message = { ...baseQueryAlliancesRequest } as QueryAlliancesRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAlliancesRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAlliancesRequest>
  ): QueryAlliancesRequest {
    const message = { ...baseQueryAlliancesRequest } as QueryAlliancesRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryAlliancesResponse: object = {};

export const QueryAlliancesResponse = {
  encode(
    message: QueryAlliancesResponse,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.assets) {
      AllianceAsset.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryAlliancesResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryAlliancesResponse } as QueryAlliancesResponse;
    message.assets = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.assets.push(AllianceAsset.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAlliancesResponse {
    const message = { ...baseQueryAlliancesResponse } as QueryAlliancesResponse;
    message.assets = [];
    if (object.assets !== undefined && object.assets !== null) {
      for (const e of object.assets) {
        message.assets.push(AllianceAsset.fromJSON(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAlliancesResponse): unknown {
    const obj: any = {};
    if (message.assets) {
      obj.assets = message.assets.map((e) =>
        e ? AllianceAsset.toJSON(e) : undefined
      );
    } else {
      obj.assets = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAlliancesResponse>
  ): QueryAlliancesResponse {
    const message = { ...baseQueryAlliancesResponse } as QueryAlliancesResponse;
    message.assets = [];
    if (object.assets !== undefined && object.assets !== null) {
      for (const e of object.assets) {
        message.assets.push(AllianceAsset.fromPartial(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryAllianceRequest: object = { denom: "" };

export const QueryAllianceRequest = {
  encode(
    message: QueryAllianceRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryAllianceRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryAllianceRequest } as QueryAllianceRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.denom = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllianceRequest {
    const message = { ...baseQueryAllianceRequest } as QueryAllianceRequest;
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = String(object.denom);
    } else {
      message.denom = "";
    }
    return message;
  },

  toJSON(message: QueryAllianceRequest): unknown {
    const obj: any = {};
    message.denom !== undefined && (obj.denom = message.denom);
    return obj;
  },

  fromPartial(object: DeepPartial<QueryAllianceRequest>): QueryAllianceRequest {
    const message = { ...baseQueryAllianceRequest } as QueryAllianceRequest;
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = object.denom;
    } else {
      message.denom = "";
    }
    return message;
  },
};

const baseQueryAllianceResponse: object = {};

export const QueryAllianceResponse = {
  encode(
    message: QueryAllianceResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.alliance !== undefined) {
      AllianceAsset.encode(message.alliance, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryAllianceResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryAllianceResponse } as QueryAllianceResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.alliance = AllianceAsset.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllianceResponse {
    const message = { ...baseQueryAllianceResponse } as QueryAllianceResponse;
    if (object.alliance !== undefined && object.alliance !== null) {
      message.alliance = AllianceAsset.fromJSON(object.alliance);
    } else {
      message.alliance = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllianceResponse): unknown {
    const obj: any = {};
    message.alliance !== undefined &&
      (obj.alliance = message.alliance
        ? AllianceAsset.toJSON(message.alliance)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllianceResponse>
  ): QueryAllianceResponse {
    const message = { ...baseQueryAllianceResponse } as QueryAllianceResponse;
    if (object.alliance !== undefined && object.alliance !== null) {
      message.alliance = AllianceAsset.fromPartial(object.alliance);
    } else {
      message.alliance = undefined;
    }
    return message;
  },
};

const baseQueryAllianceDelegationsRequest: object = { denom: "" };

export const QueryAllianceDelegationsRequest = {
  encode(
    message: QueryAllianceDelegationsRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllianceDelegationsRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllianceDelegationsRequest,
    } as QueryAllianceDelegationsRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.denom = reader.string();
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllianceDelegationsRequest {
    const message = {
      ...baseQueryAllianceDelegationsRequest,
    } as QueryAllianceDelegationsRequest;
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = String(object.denom);
    } else {
      message.denom = "";
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllianceDelegationsRequest): unknown {
    const obj: any = {};
    message.denom !== undefined && (obj.denom = message.denom);
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllianceDelegationsRequest>
  ): QueryAllianceDelegationsRequest {
    const message = {
      ...baseQueryAllianceDelegationsRequest,
    } as QueryAllianceDelegationsRequest;
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = object.denom;
    } else {
      message.denom = "";
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryAllianceDelegationsResponse: object = {};

export const QueryAllianceDelegationsResponse = {
  encode(
    message: QueryAllianceDelegationsResponse,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.delegations) {
      Delegation.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllianceDelegationsResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllianceDelegationsResponse,
    } as QueryAllianceDelegationsResponse;
    message.delegations = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.delegations.push(Delegation.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllianceDelegationsResponse {
    const message = {
      ...baseQueryAllianceDelegationsResponse,
    } as QueryAllianceDelegationsResponse;
    message.delegations = [];
    if (object.delegations !== undefined && object.delegations !== null) {
      for (const e of object.delegations) {
        message.delegations.push(Delegation.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: QueryAllianceDelegationsResponse): unknown {
    const obj: any = {};
    if (message.delegations) {
      obj.delegations = message.delegations.map((e) =>
        e ? Delegation.toJSON(e) : undefined
      );
    } else {
      obj.delegations = [];
    }
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllianceDelegationsResponse>
  ): QueryAllianceDelegationsResponse {
    const message = {
      ...baseQueryAllianceDelegationsResponse,
    } as QueryAllianceDelegationsResponse;
    message.delegations = [];
    if (object.delegations !== undefined && object.delegations !== null) {
      for (const e of object.delegations) {
        message.delegations.push(Delegation.fromPartial(e));
      }
    }
    return message;
  },
};

const baseQueryAllianceDelegationRequest: object = {
  denom: "",
  delegatorAddress: "",
};

export const QueryAllianceDelegationRequest = {
  encode(
    message: QueryAllianceDelegationRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }
    if (message.delegatorAddress !== "") {
      writer.uint32(18).string(message.delegatorAddress);
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(26).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllianceDelegationRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllianceDelegationRequest,
    } as QueryAllianceDelegationRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.denom = reader.string();
          break;
        case 2:
          message.delegatorAddress = reader.string();
          break;
        case 3:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllianceDelegationRequest {
    const message = {
      ...baseQueryAllianceDelegationRequest,
    } as QueryAllianceDelegationRequest;
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = String(object.denom);
    } else {
      message.denom = "";
    }
    if (
      object.delegatorAddress !== undefined &&
      object.delegatorAddress !== null
    ) {
      message.delegatorAddress = String(object.delegatorAddress);
    } else {
      message.delegatorAddress = "";
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllianceDelegationRequest): unknown {
    const obj: any = {};
    message.denom !== undefined && (obj.denom = message.denom);
    message.delegatorAddress !== undefined &&
      (obj.delegatorAddress = message.delegatorAddress);
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllianceDelegationRequest>
  ): QueryAllianceDelegationRequest {
    const message = {
      ...baseQueryAllianceDelegationRequest,
    } as QueryAllianceDelegationRequest;
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = object.denom;
    } else {
      message.denom = "";
    }
    if (
      object.delegatorAddress !== undefined &&
      object.delegatorAddress !== null
    ) {
      message.delegatorAddress = object.delegatorAddress;
    } else {
      message.delegatorAddress = "";
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryAllianceDelegationResponse: object = {};

export const QueryAllianceDelegationResponse = {
  encode(
    message: QueryAllianceDelegationResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.delegations !== undefined) {
      Delegation.encode(message.delegations, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllianceDelegationResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllianceDelegationResponse,
    } as QueryAllianceDelegationResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.delegations = Delegation.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllianceDelegationResponse {
    const message = {
      ...baseQueryAllianceDelegationResponse,
    } as QueryAllianceDelegationResponse;
    if (object.delegations !== undefined && object.delegations !== null) {
      message.delegations = Delegation.fromJSON(object.delegations);
    } else {
      message.delegations = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllianceDelegationResponse): unknown {
    const obj: any = {};
    message.delegations !== undefined &&
      (obj.delegations = message.delegations
        ? Delegation.toJSON(message.delegations)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllianceDelegationResponse>
  ): QueryAllianceDelegationResponse {
    const message = {
      ...baseQueryAllianceDelegationResponse,
    } as QueryAllianceDelegationResponse;
    if (object.delegations !== undefined && object.delegations !== null) {
      message.delegations = Delegation.fromPartial(object.delegations);
    } else {
      message.delegations = undefined;
    }
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Query all alliances paginated */
  Alliances(request: QueryAlliancesRequest): Promise<QueryAlliancesResponse>;
  /** Query the overall delegations of an aliance */
  Alliance(request: QueryAllianceRequest): Promise<QueryAllianceResponse>;
  /** Query the overall delegations of an aliance groupped by delegators */
  AllianceDelegations(
    request: QueryAllianceDelegationsRequest
  ): Promise<QueryAllianceDelegationsResponse>;
  /** Query the overall delegations of an aliance for a single delegator */
  AllianceDelegation(
    request: QueryAllianceDelegationRequest
  ): Promise<QueryAllianceDelegationResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("alliance.alliance.Query", "Params", data);
    return promise.then((data) => QueryParamsResponse.decode(new Reader(data)));
  }

  Alliances(request: QueryAlliancesRequest): Promise<QueryAlliancesResponse> {
    const data = QueryAlliancesRequest.encode(request).finish();
    const promise = this.rpc.request(
      "alliance.alliance.Query",
      "Alliances",
      data
    );
    return promise.then((data) =>
      QueryAlliancesResponse.decode(new Reader(data))
    );
  }

  Alliance(request: QueryAllianceRequest): Promise<QueryAllianceResponse> {
    const data = QueryAllianceRequest.encode(request).finish();
    const promise = this.rpc.request(
      "alliance.alliance.Query",
      "Alliance",
      data
    );
    return promise.then((data) =>
      QueryAllianceResponse.decode(new Reader(data))
    );
  }

  AllianceDelegations(
    request: QueryAllianceDelegationsRequest
  ): Promise<QueryAllianceDelegationsResponse> {
    const data = QueryAllianceDelegationsRequest.encode(request).finish();
    const promise = this.rpc.request(
      "alliance.alliance.Query",
      "AllianceDelegations",
      data
    );
    return promise.then((data) =>
      QueryAllianceDelegationsResponse.decode(new Reader(data))
    );
  }

  AllianceDelegation(
    request: QueryAllianceDelegationRequest
  ): Promise<QueryAllianceDelegationResponse> {
    const data = QueryAllianceDelegationRequest.encode(request).finish();
    const promise = this.rpc.request(
      "alliance.alliance.Query",
      "AllianceDelegation",
      data
    );
    return promise.then((data) =>
      QueryAllianceDelegationResponse.decode(new Reader(data))
    );
  }
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
}

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
