
      TRUNCATE acc_ap_ar;
      TRUNCATE acc_ap_ar_payment;
      TRUNCATE acc_ap_ar_payment_detail;
      TRUNCATE acc_ap_commission_payment;
      TRUNCATE acc_ap_commission_payment_detail;
      TRUNCATE acc_ap_refund_deposit_payment;
      TRUNCATE acc_ap_refund_deposit_payment_detail;
      TRUNCATE acc_cash_sale_recon;
      TRUNCATE acc_cfg_init_bank_account;
      TRUNCATE acc_close_month;
      TRUNCATE acc_close_year;
      TRUNCATE acc_credit_card_recon;
      TRUNCATE acc_credit_card_recon_detail;
      TRUNCATE acc_foreign_cash;
      TRUNCATE acc_import_journal_log;
      TRUNCATE acc_journal;
      TRUNCATE acc_journal_detail;
      TRUNCATE acc_prepaid_expense;
      TRUNCATE acc_prepaid_expense_posted;
      TRUNCATE ast_cfg_init_shipping_address;
      TRUNCATE ast_user_sub_department;
      TRUNCATE audit_log;
      TRUNCATE ban_booking;
      TRUNCATE ban_cfg_init_venue_combine;
      TRUNCATE ban_cfg_init_venue_combine_detail;
      TRUNCATE ban_cfg_init_venue_group;
      TRUNCATE ban_reservation;
      TRUNCATE ban_reservation_charge;
      TRUNCATE ban_reservation_remark;
      TRUNCATE breakfast_list_temp;
      TRUNCATE budget_expense;
      TRUNCATE budget_fb;
      TRUNCATE budget_income;
      TRUNCATE budget_statistic;
      TRUNCATE cash_count;
      /*TRUNCATE cfg_init_account
      TRUNCATE cfg_init_account_sub_group;
      TRUNCATE cfg_init_bed_type;
      TRUNCATE cfg_init_competitor_category;
      TRUNCATE cfg_init_credit_card_charge;
      TRUNCATE cfg_init_custom_lookup_field01;
      TRUNCATE cfg_init_custom_lookup_field02;
      TRUNCATE cfg_init_custom_lookup_field03;
      TRUNCATE cfg_init_custom_lookup_field04;
      TRUNCATE cfg_init_custom_lookup_field05;
      TRUNCATE cfg_init_custom_lookup_field06;
      TRUNCATE cfg_init_custom_lookup_field07;
      TRUNCATE cfg_init_custom_lookup_field08;
      TRUNCATE cfg_init_custom_lookup_field09;
      TRUNCATE cfg_init_custom_lookup_field10;
      TRUNCATE cfg_init_custom_lookup_field11;
      TRUNCATE cfg_init_custom_lookup_field12;
      TRUNCATE cfg_init_department;
      TRUNCATE cfg_init_is_fb_sub_department_group;
      TRUNCATE cfg_init_is_fb_sub_department_group_detail;
      TRUNCATE cfg_init_journal_account;
      TRUNCATE cfg_init_journal_account_category;
      TRUNCATE cfg_init_journal_account_sub_group;
      TRUNCATE cfg_init_loan_item;
      TRUNCATE cfg_init_market;
      TRUNCATE cfg_init_market_category;
      TRUNCATE cfg_init_sales;
      TRUNCATE cfg_init_purpose_of;
      TRUNCATE cfg_init_member_point_type;
      TRUNCATE cfg_init_owner;
      TRUNCATE cfg_init_package;
      TRUNCATE cfg_init_package_breakdown;
      TRUNCATE cfg_init_package_business_source;
      TRUNCATE cfg_init_printer;
      TRUNCATE cfg_init_regency;
      TRUNCATE cfg_init_reservation_mark;
      TRUNCATE cfg_init_room;
      TRUNCATE cfg_init_room_allotment_type;
      TRUNCATE cfg_init_room_amenities;
      TRUNCATE cfg_init_room_boy;
      TRUNCATE cfg_init_room_rate;
      TRUNCATE cfg_init_room_rate_breakdown;
      TRUNCATE cfg_init_room_rate_business_source;
      TRUNCATE cfg_init_room_rate_category;
      TRUNCATE cfg_init_room_rate_currency;
      TRUNCATE cfg_init_room_rate_dynamic;
      TRUNCATE cfg_init_room_rate_sub_category;
      TRUNCATE cfg_init_room_type;
      TRUNCATE cfg_init_sub_department;*/
      TRUNCATE cm_update;
      TRUNCATE company;
      TRUNCATE competitor;
      TRUNCATE competitor_data;
      TRUNCATE contact_person;
      TRUNCATE cor_cfg_init_unit;
      TRUNCATE credit_card;
      TRUNCATE data_analysis;
      TRUNCATE data_analysis_query;
      TRUNCATE data_analysis_query_list;
      TRUNCATE events;
     /* TRUNCATE fa_cfg_init_item;
      TRUNCATE fa_cfg_init_item_category;
      TRUNCATE fa_cfg_init_location;
      TRUNCATE fa_cfg_init_manufacture;*/
      TRUNCATE fa_depreciation;
      TRUNCATE fa_list;
      TRUNCATE fa_location_history;
      TRUNCATE fa_purchase_order;
      TRUNCATE fa_purchase_order_detail;
      TRUNCATE fa_receive;
      TRUNCATE fa_receive_detail;
      TRUNCATE fa_repair;
      TRUNCATE fa_revaluation;
      TRUNCATE fb_statistic;
      TRUNCATE folio;
      TRUNCATE folio_routing;
      TRUNCATE forecast_in_house_change_pax;
      TRUNCATE forecast_monthly_day;
      TRUNCATE forecast_monthly_day_previous;
      TRUNCATE grid_properties;
      TRUNCATE guest_breakdown;
      TRUNCATE guest_deposit;
      TRUNCATE guest_detail;
      TRUNCATE guest_extra_charge;
      TRUNCATE guest_extra_charge_breakdown;
      TRUNCATE guest_group;
      TRUNCATE guest_in_house;
      TRUNCATE guest_in_house_breakdown;
      TRUNCATE guest_loan_item;
      TRUNCATE guest_message;
      TRUNCATE guest_profile;
      TRUNCATE guest_scheduled_rate;
      TRUNCATE guest_to_do;
      /*TRUNCATE inv_cfg_init_item;
      TRUNCATE inv_cfg_init_item_category;
      TRUNCATE inv_cfg_init_item_category_other_cogs;
      TRUNCATE inv_cfg_init_item_category_other_cogs2;
      TRUNCATE inv_cfg_init_item_category_other_expense;
      TRUNCATE inv_cfg_init_item_uom;
      TRUNCATE inv_cfg_init_market_list;
      TRUNCATE inv_cfg_init_store;*/
      TRUNCATE inv_close_log;
      TRUNCATE inv_cost_recipe;
      TRUNCATE inv_costing;
      TRUNCATE inv_costing_detail;
      TRUNCATE inv_opname;
      TRUNCATE inv_production;
      TRUNCATE inv_purchase_order;
      TRUNCATE inv_purchase_order_detail;
      TRUNCATE inv_purchase_request;
      TRUNCATE inv_purchase_request_detail;
      TRUNCATE inv_receiving;
      TRUNCATE inv_receiving_detail;
      TRUNCATE inv_return_stock;
      TRUNCATE inv_stock_transfer;
      TRUNCATE inv_stock_transfer_detail;
      TRUNCATE inv_store_requisition;
      TRUNCATE inv_store_requisition_detail;
      TRUNCATE invoice;
      TRUNCATE invoice_item;
      TRUNCATE invoice_payment;
      TRUNCATE log;
      TRUNCATE log_backup;
      TRUNCATE log_keylock;
      TRUNCATE log_shift;
      TRUNCATE log_special_access;
      TRUNCATE log_user;
      TRUNCATE lost_and_found;
      TRUNCATE market_statistic;
      TRUNCATE member;
      TRUNCATE member_gift;
      TRUNCATE member_point;
      TRUNCATE member_point_redeem;
      TRUNCATE notif_tp;
      TRUNCATE notification;
      TRUNCATE one_time_password;
      TRUNCATE pabx_smdr;
      TRUNCATE phone_book;
      TRUNCATE pos_captain_order;
      TRUNCATE pos_captain_order_transaction;
      /*TRUNCATE pos_cfg_init_discount_limit;
      TRUNCATE pos_cfg_init_member_outlet_discount;
      TRUNCATE pos_cfg_init_member_outlet_discount_detail;
      TRUNCATE pos_cfg_init_member_product_discount;
      TRUNCATE pos_cfg_init_outlet;
      TRUNCATE pos_cfg_init_product;
      TRUNCATE pos_cfg_init_product_category;
      TRUNCATE pos_cfg_init_product_group;
      TRUNCATE pos_cfg_init_room_boy;
      TRUNCATE pos_cfg_init_table;
      TRUNCATE pos_cfg_init_table_type;
      TRUNCATE pos_cfg_init_tenan;
      TRUNCATE pos_cfg_init_waitress;*/
      TRUNCATE pos_check;
      TRUNCATE pos_check_transaction;
      TRUNCATE pos_information;
      TRUNCATE pos_iptv_menu_order;
      TRUNCATE pos_member;
      TRUNCATE pos_product_costing;
      TRUNCATE pos_reservation;
      TRUNCATE pos_reservation_table;
      TRUNCATE pos_table_unavailable;
      TRUNCATE pos_user_group_outlet;
      TRUNCATE proforma_invoice_detail;
      TRUNCATE receipt;
      TRUNCATE report_custom;
      TRUNCATE report_custom_favorite;
      TRUNCATE report_pivot_temp;
      TRUNCATE report_room_rate_structure_temp;
      TRUNCATE report_room_sales;
      TRUNCATE reservation;
      TRUNCATE reservation_extra_charge;
      TRUNCATE reservation_extra_charge_breakdown;
      TRUNCATE reservation_scheduled_rate;
      TRUNCATE room_allotment;
      TRUNCATE room_statistic;
      TRUNCATE room_status;
      TRUNCATE room_unavailable;
      TRUNCATE room_unavailable_history;
      TRUNCATE sal_activity;
      TRUNCATE sal_activity_log;
      TRUNCATE sal_notes;
      TRUNCATE sal_proposal;
      TRUNCATE sal_send_reminder;
      TRUNCATE sal_task;
      TRUNCATE sms_event;
      TRUNCATE sms_outbox;
      TRUNCATE sms_schedule;
      TRUNCATE sub_folio;
      TRUNCATE temp_sub_folio_breakdown1;
      TRUNCATE temp_sub_folio_correction_breakdown;
      TRUNCATE voucher;

      INSERT INTO audit_log (audit_date,posting_date, created_by) VALUES ("2023-06-01","2023-06-01 00:00:00","SYSTEM")