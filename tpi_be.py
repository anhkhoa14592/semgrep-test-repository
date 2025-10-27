import logging
from http import HTTPStatus
from typing import Optional, List
from urllib.parse import urljoin

import requests
from cachetools import cached, TTLCache
from cachetools.keys import hashkey
from flask import g
from injector import inject
from werkzeug.exceptions import Unauthorized

from retail_promotion import settings
from retail_promotion.const import SECRET_KEY
from retail_promotion.schema.tpi_be.best_price_guarantee import BestPriceGuaranteeSchema
from retail_promotion.schema.tpi_be.product import (
    ProductDealValidationInfo,
    ProductWithMinimumPriceCompetitor,
    ProductSKUWithSubcat,
)
from retail_promotion.settings import TPI_BE_HOST, RETAIL_DEAL_MANAGEMENT_SERVICE_APIKEY

CANNOT_GET_BGP_MESSAGE = "Error when calling to tpi be for product id {}"
CANNOT_GET_DEAL_VALIDATION_INFO = (
    "Error calling to TPI BE for deal validation info of product ID: {}"
)
CANNOT_GET_PRODUCT_BY_SKU = "Error calling to TPI BE for product of SKU: {}"
CANNOT_GET_PRODUCT_BY_SKUs = "Error calling to TPI BE for product of SKUs: {}"
CANNOT_GET_PRODUCT_BY_ID = "Error calling to TPI BE for product of ID: {}"
INSUFFICIENT_PERMISSION = "Access token with insufficient permission"


class TpiBeService:
    @inject
    def __init__(self):
        self.requests = requests
        self.logger = logging.getLogger(__name__)

    def get_best_price_guarantee_info(
        self, product_id: int
    ) -> Optional[BestPriceGuaranteeSchema]:
        get_bgp_info_by_product = urljoin(
            TPI_BE_HOST, "api/products/{}/best_price_guarantee_info"
        ).format(product_id)
        headers = {}
        if g.get("token"):
            headers["authorization"] = g.get("token")
        resp = self.requests.get(get_bgp_info_by_product, headers=headers)
        if resp.status_code == HTTPStatus.OK:
            json_response = resp.json()
            return BestPriceGuaranteeSchema.from_dict(json_response)
        else:
            self.logger.warning(CANNOT_GET_BGP_MESSAGE.format(product_id))
            return None

    @cached(
        cache=TTLCache(maxsize=128, ttl=3),
        key=lambda s, p, user_email=None: hashkey(p, user_email),
    )
    def get_deal_validation_info(
        self,
        product_id: int,
        user_email: Optional[str] = None,
    ) -> Optional[ProductDealValidationInfo]:
        get_deal_validation_info_url = urljoin(
            TPI_BE_HOST, "api/products/{}/deal_validation_info"
        ).format(product_id)
        headers = {}
        if user_email:
            headers[SECRET_KEY] = settings.SECRET_KEY
            get_deal_validation_info_url += f"?user_email={user_email}"
        elif g.get("token"):
            headers["authorization"] = g.get("token")
        resp = self.requests.get(get_deal_validation_info_url, headers=headers)
        if resp.status_code == HTTPStatus.OK:
            json_response = resp.json()
            if not json_response:
                return None
            return ProductDealValidationInfo.from_dict(json_response)
        elif resp.status_code == HTTPStatus.UNAUTHORIZED:
            raise Unauthorized(INSUFFICIENT_PERMISSION)
        else:
            self.logger.warning(CANNOT_GET_DEAL_VALIDATION_INFO.format(product_id))
            return None

    def get_product_by_sku(
        self, sku: str
    ) -> Optional[ProductWithMinimumPriceCompetitor]:
        product_by_sku_url = urljoin(TPI_BE_HOST, "api/products/sku/{}").format(sku)
        headers = {}
        if g.get("token"):
            headers["authorization"] = g.get("token")
        resp = self.requests.get(product_by_sku_url, headers=headers)
        if resp.status_code == HTTPStatus.OK:
            json_response = resp.json()
            return ProductWithMinimumPriceCompetitor.from_dict(json_response)
        elif resp.status_code == HTTPStatus.BAD_REQUEST:
            return None
        else:
            self.logger.warning(CANNOT_GET_PRODUCT_BY_SKU.format(sku))
            return None

    def get_product_by_id(
            self, product_id: int
    ) -> Optional[ProductWithMinimumPriceCompetitor]:
        product_by_id_url = urljoin(TPI_BE_HOST, "api/products/{}").format(product_id)
        headers = {}
        if g.get("token"):
            headers["authorization"] = g.get("token")
        resp = self.requests.get(product_by_id_url, headers=headers)
        if resp.status_code == HTTPStatus.OK:
            json_response = resp.json()
            return ProductWithMinimumPriceCompetitor.from_dict(json_response)
        elif resp.status_code == HTTPStatus.BAD_REQUEST:
            return None
        else:
            self.logger.warning(CANNOT_GET_PRODUCT_BY_ID.format(id))
            return None

    def get_subcat_for_skus(self, skus: List[str]) -> List[ProductSKUWithSubcat]:
        product_by_sku_url = urljoin(TPI_BE_HOST, "api/products/search")
        headers = {"authorization": RETAIL_DEAL_MANAGEMENT_SERVICE_APIKEY}

        params = {"skus": skus}
        resp = self.requests.post(
            product_by_sku_url,
            headers=headers,
            json=params,
        )
        if resp.status_code == HTTPStatus.OK:
            json_response = resp.json()
            return ProductSKUWithSubcat.Schema(many=True).load(json_response)
        elif resp.status_code == HTTPStatus.BAD_REQUEST:
            return []
        else:
            self.logger.warning(CANNOT_GET_PRODUCT_BY_SKUs.format(resp))
            return []